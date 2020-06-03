package curl_to_golang_request

import (
	"log"
	"regexp"
	"strings"
	"github.com/pkg/errors"
)

type req struct {
	h []header
	u string
}

type header struct {
	key, value string
}

var hRegex = regexp.MustCompile(`-H '(.+?): (.+?)'`)

func ParseCurlCommand(curlStr string) (*req, error) {
	spaceParts := strings.Split(curlStr, " ")
	u := spaceParts[len(spaceParts)-1]
	pp := hRegex.FindAllStringSubmatch(curlStr, -1)
	log.Println("pp", pp)
	var hh []header
	for _, p := range pp {
		log.Println("p: ", p)
		if len(p) != 3 {
			return nil, errors.Errorf("not enough parts")
		}
		h := header{key: p[1], value: p[2]}
		hh = append(hh, h)
	}
	return &req{u: strings.Replace(u, "'", "", -1), h: hh}, nil
}

/*func makeReq(r req) {
	r = http.Request{}
}*/
