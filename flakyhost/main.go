package main

// Unreliable host. Checks environment variable HITS_FAIL_COUNT and fails
// in streaks of that size. If HITS_FAIL_COUNT is 0 will never fail.

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type failingState struct {
	mtx    sync.Mutex
	toggle int64
	count  int64
	up     bool
}

// Add another hit and return if this should be a failure or not
func (fail *failingState) HitOk() bool {
	if fail.toggle == 0 {
		return true
	}

	fail.mtx.Lock()
	defer fail.mtx.Unlock()

	if fail.toggle <= fail.count {
		fail.up = !fail.up
		fail.count = 0
	}
	fail.count++
	log.Printf("up: %t count: %d next flip: %d\n", fail.up, fail.count, fail.toggle)
	return fail.up
}

// concurrency safe tracking for if a hit should fail
var failTrack failingState

func handler(w http.ResponseWriter, r *http.Request) {
	if failTrack.HitOk() {
		w.Write([]byte("OK"))
	} else {
		http.Error(w, "Pretending to be down", 400)
	}
}

func main() {
	failCount, err := strconv.ParseInt(os.Getenv("HITS_FAIL_COUNT"), 10, 0)
	if err != nil {
		log.Println(os.Getenv("HITS_FAIL_COUNT"))
		log.Println(err)
		log.Println("Host never (intentionally) failing. Set HITS_FAIL_COUNT to change that.")
	} else {
		failTrack.toggle = failCount
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("0.0.0.0:80", nil))
}
