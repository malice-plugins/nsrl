package main

const tpl = `#### NSRL Database
{{- if .Results.Found }}
 - Found :white_check_mark:
{{else}}
 - Not Found :grey_question:
{{ end -}}
`

// func printMarkDownTable(nsrl Nsrl) {
// 	fmt.Println("#### NSRL Database")
// 	if nsrl.Results.Found {
// 		fmt.Println(" - Found :white_check_mark:")
// 	} else {
// 		fmt.Println(" - Not Found :grey_question:")
// 	}
// }
