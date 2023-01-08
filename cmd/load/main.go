package main

import (
	"github.com/dezer32/parser-proxyhub.me/internal/proxyhubme"
	"os"
	"strconv"
)

func main() {
	pages := 100
	path := ""
	if os.Args[1] != "" {
		pages, _ = strconv.Atoi(os.Args[1])
	}
	if os.Args[2] != "" {
		path = os.Args[2]
	}

	proxies := proxyhubme.Parse(pages, path)
	proxies.Save("proxies")
}
