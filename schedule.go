package main

import (
	"log"
	"time"
)

type Scheduler struct {
	JobQueue       chan *Crawler
	ticker         *time.Ticker
	CrawlerConfigs []CrawlerConfig
}

func NewScheduler(crawler_configs []CrawlerConfig) *Scheduler {
	scheduler := &Scheduler{
		JobQueue:       make(chan *Crawler, 100),
		ticker:         time.NewTicker(10 * time.Minute),
		CrawlerConfigs: crawler_configs,
	}

	return scheduler
}

func (s *Scheduler) Start() {
	go func() {
		for {
			select {
			case c := <-s.JobQueue:
				c.Run(3)
			}
		}
	}()

	go func() {
		for t := range s.ticker.C {
			log.Println(t)
			for _, config := range s.CrawlerConfigs {
				s.JobQueue <- NewCrawler(config)
			}
		}
	}()
}
