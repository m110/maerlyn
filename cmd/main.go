package main

import (
	"encoding/json"
	"github.com/m110/maerlyn/checks"
	"net/http"
)

func main() {
	fetcher := checks.NewCpuFetcher()
	fetcher.Start()
	defer fetcher.Stop()

	http.HandleFunc("/cpu", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		result := fetcher.GetCpu()
		json.NewEncoder(w).Encode(result)
	})
	http.HandleFunc("/mem", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		result := fetcher.GetMem()
		json.NewEncoder(w).Encode(result)
	})
	http.ListenAndServe(":8080", nil)
}
