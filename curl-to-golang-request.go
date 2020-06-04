package curl_to_golang_request

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"
)

var hRegex = regexp.MustCompile(`-H '(.+?): (.+?)' `)

func ParseCurlCommand(curlStr string) (*req, error) {
	spaceParts := strings.Split(curlStr, " ")
	u := spaceParts[len(spaceParts)-1]
	hh, err := headersWithQuotes(curlStr)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if len(hh) == 0 {
		hh, err = headersNoQuotes(curlStr, u)
	}
	return &req{u: replaceQuote(u), h: hh}, nil
}

func headersWithQuotes(curlStr string) ([]header, error) {
	pp := hRegex.FindAllStringSubmatch(curlStr, -1)
	log.Println("pp headers", pp)
	var hh []header
	for _, p := range pp {
		log.Println("p: ", p)
		if len(p) != 3 {
			return nil, errors.Errorf("not enough parts")
		}
		h := header{key: replaceQuote(p[1]), value: replaceQuote(p[2])}
		hh = append(hh, h)
	}
	return hh, nil
}

func headersNoQuotes(curlStr string, u string) ([]header, error) {
	delim := " -H "
	parts := strings.Split(strings.Replace(curlStr, u, "", 1), delim)
	log.Println("parts in headersNoQuotes:", parts)
	parts = parts[1:]
	var hh []header
	for _, p := range parts {
		log.Println("part in headersNoQuotes:", p)
		colonParts := strings.Split(p, ": ")
		if len(colonParts) != 2 {
			return nil, errors.Errorf("not enough parts")
		}
		h := header{key: colonParts[0], value: strings.TrimSpace(colonParts[1])}
		log.Println("header in headersNoQuotes", h)
		hh = append(hh, h)
	}
	return hh, nil
}

func replaceQuote(s string) string {
	return strings.Replace(s, "'", "", -1)
}

func LogPrep() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	prepareLogrus(logrus.TraceLevel)
}

func ParseCurlCommandAndMakeReq(curlStr string, simultReqs int) error {
	r, err := ParseCurlCommand(curlStr)
	if err != nil {
		return errors.WithStack(err)
	}
	wg := sync.WaitGroup{}
	for i := 1; i<=simultReqs; i++ {
		wg.Add(1)
		go func(i int){
			defer wg.Done()
			t1 := time.Now()
			res, err := r.Do(20 * time.Second)
			l := logrus.
				WithField("time", fmt.Sprintf("%s", time.Since(t1))).
				WithField("i", i)
			if err != nil {
				l.Errorf("%+v", err)
			}
			l.Infof("status %+v, len resp: %+v", res.status, len(res.body))
		}(i)
	}
	wg.Wait()
	return nil
}

func prepareLogrus(logLevel logrus.Level) {
	customFormatter := logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	}
	customFormatter.TimestampFormat = "2006-01-02 15:04:05.999999999 -0700"
	logrus.SetFormatter(&customFormatter)
	logrus.SetReportCaller(true)
	logrus.SetLevel(logLevel)
}