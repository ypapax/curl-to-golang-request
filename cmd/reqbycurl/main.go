package main

import (
	"github.com/sirupsen/logrus"
	c "github.com/ypapax/curl-to-golang-request"
	"log"
	"os"
	"strconv"
	"strings"
)

func main()  {
	c.LogPrep()
	argsWithoutProg := os.Args[1:]
	log.Println("curl command: ", argsWithoutProg)
	simultReqsStr := os.Getenv("REQS")
	var simultReqs = 1
	if len(simultReqsStr) > 0 {
		i, err := strconv.ParseInt(simultReqsStr, 10, 64)
		if err != nil {
			logrus.Fatalf("%+v", err)
		}
		simultReqs = int(i)
	}
	if err := c.ParseCurlCommandAndMakeReq(strings.Join(argsWithoutProg, " "), simultReqs); err != nil {
		log.Println("error: ", err)
	}
}

