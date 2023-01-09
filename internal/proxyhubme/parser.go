package proxyhubme

import (
	"github.com/dezer32/proxy-checker/pkg/proxy"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
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

func Parse(pages int, path string, proxiesCh chan []proxy.Proxy, wg *sync.WaitGroup) {
	if path != "" {
		client.withPath(path)
	}

	cookieCh := make(chan *http.Cookie)
	go getData(cookieCh, proxiesCh)

	wgMain.Add(pages)
	for i := 1; i <= pages; i++ {
		log.Printf("Run load %d page.", i)
		time.Sleep(5 * time.Second)
		cookieCh <- &http.Cookie{
			Name:  "page",
			Value: strconv.Itoa(i),
		}

		log.Printf("Loaded %d page.", i)
	}

	wgMain.Wait()
	wg.Done()
}

// func getData(Cookie *http.Cookie, proxiesCh chan []Proxy) {
func getData(cookieCh chan *http.Cookie, proxiesCh chan []proxy.Proxy) {
	for cookie := range cookieCh {
		client.withCookie(cookie)

		proxies := client.getProxies()
		log.Printf("Loaded proxies (%d).", len(proxies))

		wgProxy.Add(1)
		proxiesCh <- proxies

		wgMain.Done()
	}
}
