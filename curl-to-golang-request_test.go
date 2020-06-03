package curl_to_golang_request

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	LogPrep()
	os.Exit(m.Run())
}

func TestParseCurlCommand(t *testing.T) {
	type testCase struct {
		inpCurl string
		expReq  req
	}
	cases := []testCase{
		{
			inpCurl: `curl -X 'GET' -d '' -H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3' -H 'Content-Length: 0' -H 'Content-Type: application/x-www-form-urlencoded' -H 'Sec-Fetch-User: ?1' -H 'Upgrade-Insecure-Requests: 1' -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.97 Safari/537.36' 'https://somedomain.com'`,
			expReq: req{
				u: "https://somedomain.com",
				h: []header{
					{key: "Accept", value: "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3"},
					{key: "Content-Length", value: "0"},
					{key: "Content-Type", value: "application/x-www-form-urlencoded"},
					{key: "Sec-Fetch-User", value: "?1"},
					{key: "Upgrade-Insecure-Requests", value: "1"},
					{key: "User-Agent", value: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.97 Safari/537.36"},
				},
			},
		},
		{
			inpCurl: `curl -X GET -d  -H Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3 -H Content-Length: 0 -H Content-Type: application/x-www-form-urlencoded -H Sec-Fetch-User: ?1 -H Upgrade-Insecure-Requests: 1 -H User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.97 Safari/537.36 https://somedomain.com`,
			expReq: req{
				u: "https://somedomain.com",
				h: []header{
					{key: "Accept", value: "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3"},
					{key: "Content-Length", value: "0"},
					{key: "Content-Type", value: "application/x-www-form-urlencoded"},
					{key: "Sec-Fetch-User", value: "?1"},
					{key: "Upgrade-Insecure-Requests", value: "1"},
					{key: "User-Agent", value: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.97 Safari/537.36"},
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.inpCurl, func(t *testing.T) {
			as := assert.New(t)
			a, err := ParseCurlCommand(c.inpCurl)
			if !as.NoError(err) {
				return
			}
			if !as.Equal(&c.expReq, a) {
				return
			}
		})
	}
}
