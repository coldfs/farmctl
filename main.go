package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/go-resty/resty"
)

func main() {
	usage := `farmctl
	Usage:
	farmctl env.d
`
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println(usage)
		return
	}

	resp, _ := resty.R().Get("https://deployer.service.s:31443/v1/environment/" + args[0] + "/")

	var dat map[string]interface{}

	if err := json.Unmarshal(resp.Body(), &dat); err != nil {
		panic(err)
	}

	//Будем честны, я не умею в GO. поэтому всю работу с json я буду делать по хардкору и костылями
	resources := dat["environment"].(map[string]interface{})["resources"].(map[string]interface{})
	tWriter := new(tabwriter.Writer)
	tWriter.Init(os.Stdout, 2, 8, 1, '\t', 0)
	fmt.Fprintf(tWriter, "%s\t%s\t%s\t%s\t%s\n", "NAME", "HOST", "PROVIDER", "IPS", "PORT")
	for k, v := range resources {
		data := v.(map[string]interface{})

		ips := []string{}
		if data["ips"] != nil {
			for _, ip := range data["ips"].([]interface{}) {
				ips = append(ips, ip.(string))
			}
		}
		if data["port"] == nil {
			data["port"] = ""
		}

		if data["host"] == nil {
			data["host"] = ""
		}
		if data["provider"] == nil {
			data["provider"] = ""
		}
		//

		fmt.Fprintf(tWriter, "%s\t%s\t%s\t%s\t%s\n", k, data["host"], data["provider"], strings.Join(ips, " "), data["port"])
	}
	tWriter.Flush()
}
