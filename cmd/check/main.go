package main

import (
	"encoding/json"
	"github.com/dezer32/parser-proxyhub.me/internal/proxyhubme"
	"log"
	"os"
)

func main() {
	fileName := "proxies.json"

	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	var proxies []proxyhubme.Proxy
	err = json.Unmarshal(data, &proxies)
	if err != nil {
		log.Fatal(err)
	}

	//proxies[0].Check()

	proxyCh := make(chan proxyhubme.Proxy)
	go check(proxyCh)

	for _, proxy := range proxies {
		proxyCh <- proxy
	}

	log.Printf("Done")
}

func check(proxyCh chan proxyhubme.Proxy) {
	for proxy := range proxyCh {
		proxy.Check()
	}
}
