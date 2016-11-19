package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	var page AlertManagerData
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&page)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if len(page.Alerts) == 0 {
		http.Error(w, "No alerts to display", 400)
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
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
