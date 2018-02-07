// ip-pool project main.go
package main

import (
	"io/ioutil"
	"log"

	"github.com/go-bongo/bongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/yaml.v2"
)

var config = &bongo.Config{
	ConnectionString: "localhost:27017",
	Database:         "ip_pool",
}

var forever = make(chan int)

func init() {
	conn, err := bongo.Connect(config)
	index := mgo.Index{
		Key:        []string{"ip", "port"},
		Unique:     true,
		Background: true, // See notes.
	}
	err = conn.Collection("ipdocs").Collection().EnsureIndex(index)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	config_file, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatalln("yamlFile.Get err   #%v", err)
	}
	var crawler_configs []CrawlerConfig
	err = yaml.Unmarshal([]byte(config_file), &crawler_configs)
	if err != nil {
		log.Fatalln("yamlFile.parse err   #%v", err)
	}

	c := NewChecker()
	c.Start()

	s := NewScheduler(crawler_configs)
	s.Start()

	<-forever
}
