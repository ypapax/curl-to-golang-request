package curl_to_golang_request

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"unicode/utf8"
)

type req struct {
	h []header
	u string
}

type header struct {
	key, value string
}

func (r *req) Do(timeout time.Duration) error {
	client := http.Client{Timeout: timeout}
	req, err := http.NewRequest("GET", r.u, nil)
	if err != nil {
		return errors.Wrap(err, "couldn't create request")
	}
	req.Close = true
	for _, v := range r.h {
		req.Header.Add(v.key, v.value)
		log.Println("header: ", v.key, v.value)
	}

	log.Println("requesting", r.u)
	res, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "couldn't make request, timeout: %s", timeout)
	}
	if res.StatusCode > 399 || res.StatusCode < 200 {
		var bodyText string
		defer res.Body.Close()
		b, err2 := ioutil.ReadAll(res.Body)
		if err2 != nil {
			bodyText = fmt.Sprintf("couldn't get body: %+v", err2)
		} else {
			bodyText = string(b)
		}
		const maxBodyTextChars = 2500
		var bodyTextForErr string
		if utf8.RuneCountInString(bodyText) > maxBodyTextChars {
			bodyTextForErr = string([]rune(bodyText)[:maxBodyTextChars])
		}
		err := errors.WithStack(fmt.Errorf("not good status code %+v, bodyTextForErr: %+v", res.StatusCode, bodyTextForErr))
		return err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrapf(err, "couldn't read body")
	}
	if len(b) == 0 || string(b) == "" {
		err := errors.Errorf("empty body in response status code: %+v", res.StatusCode)
		return err
	}
	log.Println("body ", string(b))
	log.Println("status ", res.StatusCode)
	return nil
}
