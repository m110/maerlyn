package main

import (
	"fmt"
	"github.com/m110/maerlyn/checks"
)

func main() {
	fetcher := checks.NewCpuFetcher()

	fetcher.Start()
	fmt.Println("CPU:", fetcher.GetTotal())
	fetcher.Stop()
}
