package main

import (
	c "github.com/ypapax/curl-to-golang-request"
	"log"
	"os"
	"strings"
)

func main()  {
	c.LogPrep()
	argsWithoutProg := os.Args[1:]
	log.Println("curl command: ", argsWithoutProg)
	if err := c.ParseCurlCommandAndMakeReq(strings.Join(argsWithoutProg, " ")); err != nil {
		log.Println("error: ", err)
	}
}

