package proxyhubme

import (
	"log"
	"net/http"
	"strconv"
	"sync"
)

const (
	baseUrl = "https://proxyhub.me/"
	//baseUrl = "https://httpbin.org/headers"
)

var (
	client  ProxyhubMe
	wgMain  sync.WaitGroup
	wgProxy sync.WaitGroup
)

func init() {
	client.init()
	wgMain = sync.WaitGroup{}
	wgProxy = sync.WaitGroup{}
}

func Parse() []Proxy {
	var res []Proxy
	cookieCh := make(chan *http.Cookie)
	proxiesCh := make(chan []Proxy)
	go getData(cookieCh, proxiesCh)

	pages := 100
	wgMain.Add(pages)
	for i := 1; i <= pages; i++ {
		log.Printf("Run load %d page.", i)
		//getData(&http.Cookie{
		//	Name:  "page",
		//	Value: strconv.Itoa(i),
		//}, proxiesCh)
		cookieCh <- &http.Cookie{
			Name:  "page",
			Value: strconv.Itoa(i),
		}

		res = append(res, <-proxiesCh...)
		log.Printf("Loaded %d page.", i)
	}

	wgMain.Wait()

	return res
}

// func getData(Cookie *http.Cookie, proxiesCh chan []Proxy) {
func getData(cookieCh chan *http.Cookie, proxiesCh chan []Proxy) {
	for cookie := range cookieCh {
		client.withCookie(cookie)

		proxies := client.getProxies()
		log.Printf("Loaded proxies (%d).", len(proxies))

		wgProxy.Add(1)
		proxiesCh <- proxies

		wgMain.Done()
	}
}
