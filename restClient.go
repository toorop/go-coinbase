package coinbase

import (
	"crypto/hmac"
	"crypto/sha256"
	//"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	//"os"
	"strings"
	"time"
)

type rClient struct {
	ak           string
	as           string
	api_endpoint string
	httpClient   *http.Client
}

func newRestClient(ak string, as string) (c *rClient) {
	return &rClient{ak, as, API_ENDPOINT, &http.Client{}}
}

// response represents usefull parts of a raw http response
type response struct {
	StatusCode int
	Status     string
	Body       []byte
}

// HandleErr return error on:
// - if err != nil
// - unexpected HTTP code
func (resp *response) HandleErr(err error, expectedHttpCode []int) error {
	if err != nil {
		return err
	}
	for _, code := range expectedHttpCode {
		if resp.StatusCode == code {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("%d - %s", resp.StatusCode, resp.Status))
}

// Do process the request & return a response (or error)
func (c *rClient) Do(method string, ressource string, payload string) (response response, err error) {
	query := fmt.Sprintf("%s/%s", c.api_endpoint, ressource)

	req, err := http.NewRequest(method, query, strings.NewReader(payload))
	if err != nil {
		return
	}
	if method == "POST" || method == "PUT" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
	}
	req.Header.Add("Accept", "application/json")

	nonce := fmt.Sprintf("%d", int64(time.Now().UnixNano()))
	mac := hmac.New(sha256.New, []byte(c.as))
	mac.Write([]byte(fmt.Sprintf("%s%s%s", nonce, query, payload)))

	req.Header.Add("ACCESS_NONCE", nonce)
	req.Header.Add("ACCESS_KEY", c.ak)
	req.Header.Add("ACCESS_SIGNATURE", fmt.Sprintf("%x", mac.Sum(nil)))
	//fmt.Println(req)

	r, err := c.httpClient.Do(req)
	if err != nil {
		return
	}
	defer r.Body.Close()

	response.StatusCode = r.StatusCode
	response.Status = r.Status
	response.Body, err = ioutil.ReadAll(r.Body)
	//fmt.Println(string(response.Body))
	return
}
