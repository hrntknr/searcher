package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/k3a/html2text"
)

func main() {
	if err := _main(); err != nil {
		panic(err)
	}
}

func _main() error {
	pUrl := flag.String("url", "", "URL")
	pHost := flag.String("host", "http://localhost:8080", "Host")
	flag.Parse()
	url := *pUrl
	host := *pHost

	if url == "" {
		return fmt.Errorf("Error: empty URL!\n")
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	plain := html2text.HTML2Text(string(html))

	body, err := json.Marshal(ReqBody{
		Uri:  url,
		Body: plain,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		host+"/regist",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	registResp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer registResp.Body.Close()

	result, err := ioutil.ReadAll(registResp.Body)
	if err != nil {
		return err
	}
	fmt.Print(string(result))
	return nil
}

type ReqBody struct {
	Uri  string `json:"uri"`
	Body string `json:"body"`
}
