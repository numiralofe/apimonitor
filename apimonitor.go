package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/consul/api"
)

type Services struct {
	Services []Service `json:"services"`
}

type Service struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	ReturnCode string `json:"return_code"`
	ReturnMsg  string `json:"return_msg"`
	Url        string `json:"url"`
}

func main() {
	// read initial paramenters
	port := flag.Int("port", 8080, "port that application should listening to, defaults to 8080")
	consulhost := flag.String("consulhost", "localhost:8500", "consul host to use, defaults to localhost:8500")
	branch := flag.String("branch", "master", "branch to use, defaults to master")
	flag.Parse()

	branch2bind := *branch
	consulhost2bind := *consulhost

	// create consul client using default config
	config := api.DefaultConfig()
	config.Address = consulhost2bind
	consul, err := api.NewClient(config)

	kv := consul.KV()
	pair, _, err := kv.Get(branch2bind+"/apimonitor.json", nil)
	if err != nil {
		panic(err)
	}

	var services Services
	var status string

	// unmarshal key pair value content from consul into 'services'
	json.Unmarshal(pair.Value, &services)

	// starting webserver handling function
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// handle curl user-agent detctection
		var useragent string
		ua := r.Header.Get("User-Agent")
		matched, err := regexp.MatchString(".*curl.*", ua)
		if (err == nil) && (matched == true) {
			useragent = "curl"
		}

		// write initial http headers
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// start priniting output
		if useragent == "curl" {
			fmt.Fprintf(w, "Service Check Status\n")
		} else {
			fmt.Fprintf(w, "<link href=\"https://maxcdn.bootstrapcdn.com/bootstrap/3.4.1/css/bootstrap.min.css\" rel=\"stylesheet\">")
			fmt.Fprintf(w, "<table class=\"table\">")
			fmt.Fprintf(w, "<thead><tr><td><b>Service Name</b></td><td><b>Response Code</b></td><td><b>Service Status</b></td><td><b>Checked URL</b></td></tr></thead><tbody>")
		}

		// iterate through every service within services array
		for i := 0; i < len(services.Services); i++ {
			// set timeout to check service availability
			client := http.Client{
				Timeout: time.Duration(1000 * time.Millisecond),
			}

			resp, err := client.Get(services.Services[i].Url)
			if (err == nil) && (strconv.Itoa(resp.StatusCode) == services.Services[i].ReturnCode) {
				status = services.Services[i].ReturnMsg
			} else {
				status = "StatusDOWN"
			}

			if useragent == "curl" {
				fmt.Fprintf(w, " "+services.Services[i].Name+" | "+strconv.Itoa(resp.StatusCode)+" | "+status+" | "+services.Services[i].Url+"\n")
			} else {
				fmt.Fprintf(w, "<tr><td>"+services.Services[i].Name+"</td><td>"+strconv.Itoa(resp.StatusCode)+"</td><td>"+status+"</td><td>"+services.Services[i].Url+"</td></tr>")
			}

		}
		if useragent != "curl" {
			fmt.Fprintf(w, "</tbody></table>")
		}
	})

	// printout used settings for the application
	fmt.Println("   ")
	fmt.Println("- listening on port", *port, "(use -port flag to set a different port)")
	fmt.Println("- connecting to consul on", *consulhost, "(use -consul flag to set a different consul host)")
	fmt.Println("- using", *branch, "branch (use -branch flag to set a different branch)")

	// bind to port and start http listener
	port2bind := strconv.Itoa(*port)
	log.Fatal(http.ListenAndServe(":"+port2bind, nil))

}
