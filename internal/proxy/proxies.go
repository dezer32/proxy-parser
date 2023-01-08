package proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type ByCountry map[string][]Proxy

type Proxies struct {
	List    []Proxy
	Grouped ByCountry
}

func (p *Proxies) GroupingByCountry() {
	p.Grouped = make(map[string][]Proxy)

	for _, proxy := range p.List {
		city := strings.ToLower(proxy.City)
		p.Grouped[city] = append(p.Grouped[city], proxy)
	}
}

func (p *Proxies) SaveGroup(prefixFileName string) {
	p.GroupingByCountry()
	for c, i := range p.Grouped {
		fileName := fmt.Sprintf("%s.%s.json", prefixFileName, c)
		p.save(fileName, i)
	}

}

func (p *Proxies) Save(prefixFileName string) {
	fileName := fmt.Sprintf("%s.full.json", prefixFileName)
	p.save(fileName, p.List)
}

func (p *Proxies) save(fileName string, list []Proxy) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Open file failed. Error: %s", err)
	}
	defer file.Close()

	j, _ := json.Marshal(list)

	log.Printf("Save proxies (%d) to %s.", len(list), fileName)

	file.Write(j)
}

func (p *Proxies) Load(fileName string) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &p.List)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Load proxies (%d).", len(p.List))
}

func (p *Proxies) Remove(i int) {
	p.List[i] = p.List[len(p.List)-1]
	p.List = p.List[:len(p.List)-1]
}
