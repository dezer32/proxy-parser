package main

import (
	"encoding/json"
	"github.com/dezer32/parser-proxyhub.me/internal/proxyhubme"
	"log"
	"os"
)

func main() {
	file, err := os.Create("proxies.json")
	if err != nil {
		log.Fatalf("Open file failed. Error: %s", err)
	}
	defer file.Close()

	proxies := proxyhubme.Parse()
	j, _ := json.Marshal(proxies)

	log.Printf("Done loaded proxies (%d).", len(proxies))

	file.Write(j)
}
