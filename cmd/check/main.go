package main

import (
	"github.com/dezer32/parser-proxyhub.me/internal/proxy"
	"log"
	"os"
	"sync"
	"time"
)

var (
	wgMain  = sync.WaitGroup{}
	wgProxy = sync.WaitGroup{}
)

func main() {
	fileName := "proxies.json"
	if os.Args[1] != "" {
		fileName = os.Args[1]
	}
	proxies := proxy.Proxies{}
	proxies.Load(fileName)

	verifiedProxyChan := make(chan proxy.Proxy)
	for _, p := range proxies.List {
		wgMain.Add(1)
		go func(rawProxy proxy.Proxy, verifiedProxyCh chan proxy.Proxy) {
			rawProxy.Check()
			wgMain.Done()

			if rawProxy.IsWorked == true {
				wgProxy.Add(1)
				verifiedProxyChan <- rawProxy

			}
		}(p, verifiedProxyChan)
	}

	wgMain.Wait()

	verifiedProxies := proxy.Proxies{}
	timeout := time.After(15 * time.Second)

	isBreak := false
	for !isBreak {
		select {
		case verifiedProxy := <-verifiedProxyChan:
			verifiedProxies.List = append(verifiedProxies.List, verifiedProxy)
			verifiedProxies.Save("temp.checked.proxies")
			wgProxy.Done()
		case <-timeout:
			isBreak = true
		}
	}

	verifiedProxies.Save("checked")

	wgProxy.Wait()

	log.Printf("Done")
}

//func check(proxyCh chan *proxy.Proxy) {
//	for p := range proxyCh {
//		p.Check()
//		wg.Done()
//	}
//}
