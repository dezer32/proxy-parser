package proxyhubme

import (
	"errors"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strconv"
)

type ProxyhubMe struct {
	Client  http.Client
	Cookie  *http.Cookie
	Headers map[string]string
}

func (p *ProxyhubMe) init() {
	jar, err := cookiejar.New(nil)
	logErr(err)

	p.Client = http.Client{
		Jar: jar,
	}

	p.Headers = map[string]string{
		"authority":                 "proxyhub.me",
		"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8",
		"accept-language":           "ru-RU,ru;q=0.9",
		"sec-fetch-dest":            "document",
		"sec-fetch-mode":            "navigate",
		"sec-fetch-site":            "none",
		"sec-fetch-user":            "?1",
		"sec-gpc":                   "1",
		"upgrade-insecure-requests": "1",
		"user-agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
	}
}

func (p *ProxyhubMe) withCookie(cookie *http.Cookie) {
	p.Cookie = cookie
}

func (p *ProxyhubMe) getBody() io.ReadCloser {
	req, err := http.NewRequest("GET", baseUrl, nil)
	logErr(err)
	req.AddCookie(p.Cookie)

	for name, value := range p.Headers {
		req.Header.Set(name, value)
	}

	resp, err := p.Client.Do(req)
	logErr(err)
	log.Printf("Loaded page (%s).", p.Cookie.Value)

	return resp.Body
}

func (p *ProxyhubMe) getProxies() []Proxy {
	body := p.getBody()
	defer body.Close()

	doc, err := html.Parse(p.getBody())
	logErr(err)
	tn, err := getFragment(doc, "tbody")
	logErr(err)

	var res []Proxy
	for child := tn.FirstChild; child != nil; child = child.NextSibling {
		res = append(res, parseProxy(child))
	}

	log.Printf("Parsed proxies on page (%s).", p.Cookie.Value)

	return res
}

func parseProxy(doc *html.Node) Proxy {
	iter := 0
	p := Proxy{}

	for child := doc.FirstChild; child != nil; child = child.NextSibling {
		iter++
		switch iter {
		case 1:
			p.Ip = child.FirstChild.Data
		case 2:
			p.Port, _ = strconv.Atoi(child.FirstChild.Data)
		case 3:
			p.Protocol = child.FirstChild.Data
		case 4:
			p.Anonymity = child.FirstChild.Data
		case 5:
			p.Country = child.LastChild.FirstChild.Data
		case 6:
			p.City = child.FirstChild.Data
		}
	}

	return p
}

func getFragment(doc *html.Node, tagName string) (*html.Node, error) {
	var body *html.Node
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == tagName {
			body = node
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	if body != nil {
		return body, nil
	}
	return nil, errors.New("Missing <body> in the doc tree")
}

func logErr(err error) {
	if err != nil {
		log.Fatalf("Error occured. Error is: %s", err.Error())
	}
}
