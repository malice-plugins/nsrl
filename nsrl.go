package main

import (
	"bytes"
	"encoding/binary"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/fatih/structs"
	"github.com/gorilla/mux"
	"github.com/malice-plugins/pkgs/database"
	"github.com/malice-plugins/pkgs/database/elasticsearch"
	"github.com/malice-plugins/pkgs/utils"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"github.com/willf/bloom"
)

var (
	// Version stores the plugin's version
	Version string

	// BuildTime stores the plugin's build time
	BuildTime string

	// ErrorRate stores the bloomfilter desired error-rate
	ErrorRate string

	// HashType is the type of hash to use to build the bloomfilter
	HashType string

	// es is the elasticsearch database object
	es elasticsearch.Database
)

const (
	// NSRL fields
	sha1         = 0
	md5          = 1
	crc32        = 2
	fileName     = 3
	fileSize     = 4
	productCode  = 5
	opSystemCode = 6
	specialCode  = 7
)

const (
	name     = "nsrl"
	category = "intel"
)

type pluginResults struct {
	ID   string      `json:"id" gorethink:"id,omitempty"`
	Data ResultsData `json:"nsrl" gorethink:"nsrl"`
}

// Nsrl json object
type Nsrl struct {
	Results ResultsData `json:"nsrl"`
}

// ResultsData json object
type ResultsData struct {
	Found    bool   `json:"found"`
	Hash     string `json:"hash"`
	MarkDown string `json:"markdown,omitempty" structs:"markdown,omitempty"`
}

func assert(err error) {
	if err != nil {
		log.WithFields(log.Fields{
			"plugin":   name,
			"category": category,
			"path":     path,
		}).Fatal(err)
	}
}

func generateMarkDownTable(n Nsrl) string {
	var tplOut bytes.Buffer

	t := template.Must(template.New("nsrl").Parse(tpl))

	err := t.Execute(&tplOut, n)
	if err != nil {
		log.Println("executing template:", err)
	}

	return tplOut.String()
}

func lineCounter(r io.Reader) (uint64, error) {
	buf := make([]byte, 32*1024)
	var count uint64
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += uint64(bytes.Count(buf[:c], lineSep))

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func getNSRLFieldFromHashType() int {
	switch strings.ToLower(HashType) {
	case "sha1":
		return 0
	case "md5":
		return 1
	case "crc32":
		return 2
	case "filename":
		return 3
	case "filesize":
		return 4
	case "productcode":
		return 5
	case "opsystemcode":
		return 6
	case "specialcode":
		return 7
	default:
		log.Fatal(fmt.Errorf("hash type %s not supported", HashType))
	}
	return -1
}

// build bloomfilter from NSRL database
func buildFilter() {
	var err error
	nsrlField := getNSRLFieldFromHashType()

	// open NSRL database
	nsrlDB, err := os.Open("NSRLFile.txt")
	assert(err)
	defer nsrlDB.Close()
	// count lines in NSRL database
	lines, err := lineCounter(nsrlDB)
	assert(err)
	log.Debugf("Number of lines in NSRLFile.txt: %d", lines)
	// write line count to file LINECOUNT
	buf := new(bytes.Buffer)
	assert(binary.Write(buf, binary.LittleEndian, lines))
	assert(ioutil.WriteFile("LINECOUNT", buf.Bytes(), 0644))

	// Create new bloomfilter with size = number of lines in NSRL database
	erate, err := strconv.ParseFloat(ErrorRate, 64)
	assert(err)

	filter := bloom.NewWithEstimates(uint(lines), erate)

	// jump back to the begining of the file
	_, err = nsrlDB.Seek(0, io.SeekStart)
	assert(err)

	log.Debug("Loading NSRL database into bloomfilter")
	reader := csv.NewReader(nsrlDB)
	// strip off csv header
	_, _ = reader.Read()
	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}
		assert(err)

		// log.Debug(record)
		filter.Add([]byte(record[nsrlField]))
	}

	bloomFile, err := os.Create("nsrl.bloom")
	assert(err)
	defer bloomFile.Close()

	log.Debug("Writing bloomfilter to disk")
	_, err = filter.WriteTo(bloomFile)
	assert(err)
}

// lookUp queries the NSRL bloomfilter for a hash
func lookUp(hash string, timeout int) ResultsData {

	var lines uint64
	nsrlResults := ResultsData{}

	// read line count from file LINECOUNT
	lineCount, err := ioutil.ReadFile("LINECOUNT")
	assert(err)
	buf := bytes.NewReader(lineCount)
	assert(binary.Read(buf, binary.LittleEndian, &lines))
	log.Debugf("Number of lines in NSRLFile.txt: %d", lines)

	// Create new bloomfilter with size = number of lines in NSRL database
	erate, err := strconv.ParseFloat(ErrorRate, 64)
	assert(err)

	filter := bloom.NewWithEstimates(uint(lines), erate)

	// load NSRL bloomfilter from file
	f, err := os.Open("nsrl.bloom")
	assert(err)
	_, err = filter.ReadFrom(f)
	assert(err)

	// test of existance of hash in bloomfilter
	nsrlResults.Found = filter.TestString(hash)

	return nsrlResults
}

func webService() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/lookup/{hash}", webLookUp)
	log.Info("web service listening on port :3993")
	log.Fatal(http.ListenAndServe(":3993", router))
}

func webLookUp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	if utils.StringInSlice(HashType, []string{"md5", "sha1"}) {
		hashType, _ := utils.GetHashType(hash)
		if !strings.EqualFold(hashType, strings.ToUpper(HashType)) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Please supply a proper %s hash to query", strings.ToUpper(HashType))
			return
		}
	}

	nsrl := Nsrl{Results: lookUp(strings.ToUpper(hash), 10)}
	nsrl.Results.Hash = hash

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if nsrl.Results.Found {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

	if err := json.NewEncoder(w).Encode(nsrl); err != nil {
		panic(err)
	}
}

func printStatus(resp gorequest.Response, body string, errs []error) {
	fmt.Println(body)
}

func main() {

	cli.AppHelpTemplate = utils.AppHelpTemplate
	app := cli.NewApp()

	app.Name = "nsrl"
	app.Author = "blacktop"
	app.Email = "https://github.com/blacktop"
	app.Version = Version + ", BuildTime: " + BuildTime
	app.Compiled, _ = time.Parse("20060102", BuildTime)
	app.Usage = "Malice NSRL Plugin"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, V",
			Usage: "verbose output",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "web",
			Aliases: []string{"w"},
			Usage:   "Create a NSRL lookup web service",
			Action: func(c *cli.Context) error {
				webService()
				return nil
			},
		},
		{
			Name:    "build",
			Aliases: []string{"b"},
			Usage:   "Build bloomfilter from NSRL database",
			Action: func(c *cli.Context) error {
				if c.GlobalBool("verbose") {
					log.SetLevel(log.DebugLevel)
				}

				buildFilter()
				return nil
			},
		},
		{
			Name:      "lookup",
			Aliases:   []string{"l"},
			Usage:     "Query NSRL for hash",
			ArgsUsage: fmt.Sprintf("%s to query NSRL with", strings.ToUpper(HashType)),
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "elasticsearch",
					Value:       "",
					Usage:       "elasticsearch url for Malice to store results",
					EnvVar:      "MALICE_ELASTICSEARCH_URL",
					Destination: &es.URL,
				},
				cli.BoolFlag{
					Name:   "post, p",
					Usage:  "POST results to Malice webhook",
					EnvVar: "MALICE_ENDPOINT",
				},
				cli.BoolFlag{
					Name:   "proxy, x",
					Usage:  "proxy settings for Malice webhook endpoint",
					EnvVar: "MALICE_PROXY",
				},
				cli.IntFlag{
					Name:   "timeout",
					Value:  10,
					Usage:  "malice plugin timeout (in seconds)",
					EnvVar: "MALICE_TIMEOUT",
				},
				cli.BoolFlag{
					Name:  "table, t",
					Usage: "output as Markdown table",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Args().Present() {
					hash := strings.ToUpper(c.Args().First())

					if utils.StringInSlice(HashType, []string{"md5", "sha1"}) {
						hashType, _ := utils.GetHashType(hash)

						if !strings.EqualFold(hashType, strings.ToUpper(HashType)) {
							log.Fatal(fmt.Errorf("please supply a valid %s hash to query NSRL with", strings.ToUpper(HashType)))
						}
					}

					if c.GlobalBool("verbose") {
						log.SetLevel(log.DebugLevel)
					}

					nsrl := Nsrl{Results: lookUp(hash, c.GlobalInt("timeout"))}
					nsrl.Results.Hash = hash
					nsrl.Results.MarkDown = generateMarkDownTable(nsrl)

					// upsert into Database
					if len(c.String("elasticsearch")) > 0 {
						err := es.Init()
						if err != nil {
							return errors.Wrap(err, "failed to initalize elasticsearch")
						}
						err = es.StorePluginResults(database.PluginResults{
							ID:       utils.Getopt("MALICE_SCANID", hash),
							Name:     name,
							Category: category,
							Data:     structs.Map(nsrl.Results),
						})
						if err != nil {
							return errors.Wrapf(err, "failed to index malice/%s results", name)
						}
					}

					if c.Bool("table") {
						fmt.Println(nsrl.Results.MarkDown)
					} else {
						nsrl.Results.MarkDown = ""
						nsrlJSON, err := json.Marshal(nsrl)
						assert(err)
						if c.Bool("post") {
							request := gorequest.New()
							if c.Bool("proxy") {
								request = gorequest.New().Proxy(os.Getenv("MALICE_PROXY"))
							}
							request.Post(os.Getenv("MALICE_ENDPOINT")).
								Set("X-Malice-ID", utils.Getopt("MALICE_SCANID", hash)).
								Send(string(nsrlJSON)).
								End(printStatus)

							return nil
						}
						fmt.Println(string(nsrlJSON))
					}
				} else {
					log.Fatal(fmt.Errorf("please supply a %s hash to query NSRL with", strings.ToUpper(HashType)))
				}
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	assert(err)
}
