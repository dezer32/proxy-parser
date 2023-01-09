package main

import (
	"flag"
	"fmt"
	"github.com/dezer32/parser-proxyhub.me/internal/proxyhubme"
	"github.com/dezer32/proxy-checker/pkg/proxy"
	"sync"
	"time"
)

var (
	mutex = sync.Mutex{}
	wg    = sync.WaitGroup{}

	pages    int
	page     string
	proxies  proxy.Proxies
	fileName string
)

func init() {
	wg = sync.WaitGroup{}

	defaultOutputFileName := fmt.Sprintf("proxies.%d.json", time.Now().Unix())

	flag.StringVar(&fileName, "o", defaultOutputFileName, "File to save proxies")
	flag.IntVar(&pages, "n", 100, "# pages to load")
	flag.StringVar(&page, "p", "", "different page to load")
	flag.Parse()
}

func main() {
	proxiesCh := make(chan []proxy.Proxy)
	wg.Add(1)
	go proxyhubme.Parse(pages, page, proxiesCh, &wg)
	go consumeProxies(proxiesCh)

	wg.Wait()

	proxies.Save(fileName)
}

func consumeProxies(proxiesCh chan []proxy.Proxy) {
	for p := range proxiesCh {
		mutex.Lock()
		proxies.List = append(proxies.List, p...)
		mutex.Unlock()
	}
}
