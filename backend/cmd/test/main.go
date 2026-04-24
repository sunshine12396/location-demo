package main

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/oschwald/geoip2-golang"
)

func main() {
	db, _ := geoip2.Open("resources/GeoLite2-City.mmdb")
	defer db.Close()

	ip := net.ParseIP("118.69.178.91")
	record, _ := db.City(ip)
	record, err := db.City(ip)
	if err != nil {
		fmt.Printf("%v+", err)
		return
	}
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}
