package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Checker struct {
	check_site string
	ticker     *time.Ticker
}

func NewChecker() *Checker {
	return &Checker{
		check_site: "http://www.baidu.com",
		ticker:     time.NewTicker(1 * time.Minute),
	}
}

func (c *Checker) Check() {
	results := AvaiableIP()
	var wg sync.WaitGroup

	for _, ip_doc := range results {
		wg.Add(1)
		fmt.Printf("%s:%s\n", ip_doc.Ip, ip_doc.Port)
		go func(ip, port string) {
			client := newClient(ip, port)
			_, err := client.Get(c.check_site)
			if err != nil {
				log.Println(fmt.Sprintf("error: %s:%s -- %v", ip, port, err))
				Unavailable(ip, port)
			}
			wg.Done()
		}(ip_doc.Ip, ip_doc.Port)
	}

	wg.Wait()
	fmt.Println("check done")
}

func (c *Checker) Start() {
	go func() {
		for t := range c.ticker.C {
			log.Println(t)
			c.Check()
		}
	}()
}

func newClient(ip, port string) *http.Client {
	proxy, _ := url.Parse(fmt.Sprintf("http://%s:%s", ip, port))

	transport := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}
}
