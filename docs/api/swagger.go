package api

import (
	"io/ioutil"
	"log"
)

func GetSwaggerJSON() string {
	data, err := ioutil.ReadFile("swagger.json")
	if err != nil {
		log.Fatalf("can not read file swagger.json error: %v", err)
	}
	return string(data)
}
