package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/healthcheck", nil)
	if err != nil {
		fmt.Println("creating request failed")
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("healthcheck failed")
		os.Exit(1)
	}

	if res.StatusCode != http.StatusOK {
		fmt.Println("healthcheck failed")
		os.Exit(1)
	}
}
