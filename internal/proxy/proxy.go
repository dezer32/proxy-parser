package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Proxy struct {
	Ip        string `json:"ip"`
	Port      int    `json:"port"`
	Protocol  string `json:"protocol"`
	Anonymity string `json:"anonymity"`
	Country   string `json:"country"`
	City      string `json:"city"`
	IsWorked  bool   `json:"is_checked"`
}

//type httpbinip struct {
//	Origin string `json:"origin"`
//}

func (p *Proxy) Check() bool {
	proxyUrl, _ := url.Parse(fmt.Sprintf("%s://%s:%d", strings.ToLower(p.Protocol), p.Ip, p.Port))

	client := &http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyURL(proxyUrl),
			IdleConnTimeout: 15 * time.Second,
		},
	}

	log.Printf("Check proxy %s.", proxyUrl)
	//resp, err := client.Get("https://httpbin.org/ip")
	resp, err := client.Get("https://ya.ru")
	if err != nil {
		log.Printf("For proxy %s, error %s", proxyUrl, err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Status code is %d", resp.StatusCode)
	}

	//respData, _ := io.ReadAll(resp.Body)
	//ip := httpbinip{}
	//json.Unmarshal(respData, &ip)

	p.IsWorked = resp.StatusCode == 200

	log.Printf("Proxy: %s, check: %v", proxyUrl, p.IsWorked)
	return true
}
