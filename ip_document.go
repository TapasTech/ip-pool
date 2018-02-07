package main

import (
	"github.com/go-bongo/bongo"
	//"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/bson"
)

type IPDoc struct {
	bongo.DocumentBase `bson:",inline"`
	Ip                 string
	Port               string
	Available          bool
	ResponseTime       string
	Website            string
}

func (i *IPDoc) Save() error {
	var DBConn, _ = bongo.Connect(config)
	err := DBConn.Collection("ipdocs").Save(i)
	return err
}

func AvaiableIP() []IPDoc {
	connection, _ := bongo.Connect(config)
	results := connection.Collection("ipdocs").Find(bson.M{"available": true})
	var ip_docs []IPDoc
	results.Query.All(&ip_docs)
	return ip_docs
}

func Unavailable(ip, port string) error {
	var DBConn, _ = bongo.Connect(config)
	doc := &IPDoc{}
	DBConn.Collection("ipdocs").FindOne(bson.M{"ip": ip, "port": port}, doc)
	doc.Available = false
	err := DBConn.Collection("ipdocs").Save(doc)
	return err
}
