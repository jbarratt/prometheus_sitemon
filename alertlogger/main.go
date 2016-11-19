package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request from %s: %s %s", r.RemoteAddr, r.Method, r.URL.String())
	var page AlertManagerData
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	bodyText, err := ioutil.ReadAll(r.Body)
	if err == nil {
		log.Printf("%s\n", bodyText)
	} else {
		log.Println("Error reading body")
		return
	}
	err = json.Unmarshal(bodyText, &page)
	if err != nil {
		http.Error(w, err.Error(), 400)
		log.Printf("JSON Decode Error: %v\n", err.Error())
		return
	}
	if len(page.Alerts) == 0 {
		http.Error(w, "No alerts to display", 400)
		log.Println("No alerts to display")
		return
	}
	log.Println(page.Alerts[0].Annotations.Description)
	log.Println("Status: " + page.Status)
	w.Write([]byte("OK"))
}

func main() {
	/// Set up logging to output file
	logfile := os.Getenv("LOG_ALERT_PATH")
	if logfile != "" {
		f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			fmt.Printf("Error opening log file: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("0.0.0.0:8088", nil))
}
