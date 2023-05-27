package main

import (
	"github.com/ImOlli/go-lcu/proxy"
	"log"
)

func main() {
	// Automatically resolves the port and auth token of the lcu and creates a reverse proxy
	p, err := proxy.CreateProxy(":8080")

	if err != nil {
		panic(err)
	}

	// Optionally features
	p.DisableCORS = true
	p.DisableCertCheck = true
	p.DisableStartUpMessage = true

	log.Fatal(p.ListenAndServe())
}
