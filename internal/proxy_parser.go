package internal

import (
	"github.com/dezer32/proxy-checker/pkg/proxy"
	"sync"
)

type ParserInterface interface {
	Parse(pages int, path string, proxiesCh chan []proxy.Proxy, wg *sync.WaitGroup)
}
