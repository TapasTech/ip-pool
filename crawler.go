package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mitchellh/mapstructure"

	"time"

	"github.com/antchfx/htmlquery"
)

type CrawlerConfig struct {
	Name     string            `yaml:"name"`
	PathFmt  string            `yaml:"path_fmt"`
	Start    int               `yaml:"start"`
	Xpath    string            `yaml:"xpath"`
	Position map[string]string `yaml:"position"`
}

type Crawler struct {
	Config CrawlerConfig
}

func NewCrawler(config CrawlerConfig) *Crawler {
	return &Crawler{
		Config: config,
	}

}

func (c *Crawler) Run(end int) {
	for j := 1; j <= end; j++ {
		url := fmt.Sprintf(c.Config.PathFmt, j)
		c.Spider(url)
		time.Sleep(1000 * time.Millisecond)
	}

	fmt.Println("Done")
}

func (c *Crawler) Spider(url string) {

	resp, err := http.Get(url)

	if err != nil {
		log.Fatalln("http get url error   #%v", err)
	}

	defer resp.Body.Close()

	doc, err := htmlquery.Parse(resp.Body)

	if err != nil {
		fmt.Println(err)
	}

	for _, n := range htmlquery.Find(doc, c.Config.Xpath) {
		result := make(map[string]interface{})
		for k, v := range c.Config.Position {
			vv := htmlquery.FindOne(n, v)
			result[k] = vv.FirstChild.Data
		}
		result["Available"] = true
		result["Website"] = c.Config.Name
		var ip_doc IPDoc
		mapstructure.Decode(result, &ip_doc)
		fmt.Println(ip_doc.Ip + "---" + url)
		ip_doc.Save()
	}
}
