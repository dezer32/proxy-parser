package proxyhubme

import (
	"fmt"
	"io"
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
}

func (p *Proxy) Check() bool {
	proxy := fmt.Sprintf("%s://%s:%d", strings.ToLower(p.Protocol), p.Ip, p.Port)
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: func(r *http.Request) (*url.URL, error) {
				pUrl := &url.URL{}
				pUrl.Parse(proxy)

				return pUrl, nil
			},
			IdleConnTimeout: 60 * time.Second,
		},
	}

	resp, err := client.Get("https://httpbin.org/ip")
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Status code is %d", resp.StatusCode)
	}

	respData, _ := io.ReadAll(resp.Body)

	log.Printf("Proxy: %s, status code: %d, response: %s", proxy, resp.StatusCode, respData)
	return true
}
