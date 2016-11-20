package main

import (
	"bytes"
	"encoding/binary"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/fatih/structs"
	"github.com/maliceio/go-plugin-utils/database/elasticsearch"
	"github.com/maliceio/go-plugin-utils/utils"
	"github.com/parnurzeal/gorequest"
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
	Found bool `json:"found"`
}

func printMarkDownTable(nsrl Nsrl) {
	fmt.Println("#### NSRL")
	if nsrl.Results.Found {
		fmt.Println(" - Found")
	} else {
		fmt.Println(" - Not Found")
	}
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

// build bloomfilter from NSRL database
func buildFilter() {

	// open NSRL database
	nsrlDB, err := os.Open("/nsrl/NSRLFile.txt")
	utils.Assert(err)
	// count lines in NSRL database
	lines, err := lineCounter(nsrlDB)
	log.Debugf("Number of lines in NSRLFile.txt: %d\n", lines)
	nsrlDB.Close()
	// write line count to file LINECOUNT
	buf := new(bytes.Buffer)
	utils.Assert(binary.Write(buf, binary.LittleEndian, lines))
	utils.Assert(ioutil.WriteFile("/nsrl/LINECOUNT", buf.Bytes(), 0644))

	// Create new bloomfilter with size = number of lines in NSRL database
	erate, err := strconv.ParseFloat(ErrorRate, 64)
	filter := bloom.NewWithEstimates(uint(lines), erate)

	// open NSRL database
	nsrlDB, err = os.Open("/nsrl/NSRLFile.txt")
	utils.Assert(err)
	defer nsrlDB.Close()

	reader := csv.NewReader(nsrlDB)
	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Add SHA256
		log.Debug(record)
		filter.Add([]byte(record[4]))
	}

	bloomFile, err := os.Create("/nsrl.bloom")
	utils.Assert(err)
	defer bloomFile.Close()

	filter.WriteTo(bloomFile)
}

// lookUp queries the NSRL bloomfilter for a hash
func lookUp(hash string, timeout int) ResultsData {

	var lines int
	nsrlResults := ResultsData{}

	// read line count from file LINECOUNT
	lineCount, err := ioutil.ReadFile("/nsrl/LINECOUNT")
	utils.Assert(err)
	buf := bytes.NewReader(lineCount)
	utils.Assert(binary.Read(buf, binary.LittleEndian, &lines))
	log.Debugf("Number of lines in NSRLFile.txt: %d\n", lines)

	// Create new bloomfilter with size = number of lines in NSRL database
	erate, err := strconv.ParseFloat(ErrorRate, 64)
	filter := bloom.NewWithEstimates(uint(lines), erate)

	// load NSRL bloomfilter from file
	f, err := os.Open("/nsrl.bloom")
	utils.Assert(err)
	_, err = filter.ReadFrom(f)
	utils.Assert(err)

	// test of existance of hash in bloomfilter
	nsrlResults.Found = filter.TestString(hash)

	return nsrlResults
}

func printStatus(resp gorequest.Response, body string, errs []error) {
	fmt.Println(body)
}

func main() {

	var elastic string

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
		cli.StringFlag{
			Name:        "elasitcsearch",
			Value:       "",
			Usage:       "elasitcsearch address for Malice to store results",
			EnvVar:      "MALICE_ELASTICSEARCH",
			Destination: &elastic,
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
		cli.BoolFlag{
			Name:  "table, t",
			Usage: "output as Markdown table",
		},
		cli.IntFlag{
			Name:   "timeout",
			Value:  60,
			Usage:  "malice plugin timeout (in seconds)",
			EnvVar: "MALICE_TIMEOUT",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "build",
			Aliases: []string{"b"},
			Usage:   "Build bloomfilter from NSRL database",
			Action: func(c *cli.Context) error {
				if c.GlobalBool("verbose") {
					log.SetLevel(log.DebugLevel)
				}

				// build bloomfilter
				buildFilter()
				return nil
			},
		},
		{
			Name:      "lookup",
			Aliases:   []string{"l"},
			Usage:     "Query NSRL for hash",
			ArgsUsage: "HASH to query NSRL with",
			Action: func(c *cli.Context) error {
				if c.Args().Present() {
					hash := c.Args().First()

					if c.GlobalBool("verbose") {
						log.SetLevel(log.DebugLevel)
					}

					nsrl := Nsrl{Results: lookUp(hash, c.Int("timeout"))}

					// upsert into Database
					elasticsearch.InitElasticSearch(elastic)
					elasticsearch.WritePluginResultsToDatabase(elasticsearch.PluginResults{
						ID:       utils.Getopt("MALICE_SCANID", hash),
						Name:     name,
						Category: category,
						Data:     structs.Map(nsrl.Results),
					})

					if c.GlobalBool("table") {
						printMarkDownTable(nsrl)
					} else {
						nsrlJSON, err := json.Marshal(nsrl)
						utils.Assert(err)
						if c.GlobalBool("post") {
							request := gorequest.New()
							if c.GlobalBool("proxy") {
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
					log.Fatal(fmt.Errorf("Please supply a MD5/SHA1/SHA256 hash to query NSRL with."))
				}
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	utils.Assert(err)
}
