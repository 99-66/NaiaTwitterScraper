package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const METHOD = "GET"
const URL = "https://api.twitter.com/2/tweets/sample/stream"
const TOKEN = ""

type Tweet struct {
	Id string `json:"id"`
	Text string `json:"text"`
}

type StreamTweet struct {
	Tweet `json:"data"`
}

type Header struct {
	Value string `json:"value"`
	Tag string `json:"tag"`
}

type FilterStreamPayload struct {
	Add []Header `json:"add"`
}

func setRules(bearerToken string) error {
	apiUrl := "https://api.twitter.com/2/tweets/search/stream/rules"
	method := "POST"
	payload:= FilterStreamPayload{
		Add: []Header{
			{Value: "from:twitterdev from:twitterapi has:links", Tag: "Filtering Korean Twitter"},
		},
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, apiUrl, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("rule set new request failed. %v\n", err)
	}

	req.Header.Set("Authorization", bearerToken)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("rule set running failed. %v\n", err)
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(data))

	if res.StatusCode == 200 {
		return nil
	} else {
		return fmt.Errorf("Rule Status: %v\n", err)
	}
}

func main() {
	client := &http.Client{}
	req, err := http.NewRequest(METHOD, URL, nil)
	if err != nil {
		log.Fatalf("http new request failed. %v\n", err)
	}

	bearerToken := "Bearer " + TOKEN
	req.Header.Add("Authorization", bearerToken)
	err = setRules(bearerToken)
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("client running failed. %v\n", err)
	}
	defer res.Body.Close()

	reader := bufio.NewReader(res.Body)
	for {
		line, err := reader.ReadBytes('\r')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		var tweet StreamTweet
		err = json.Unmarshal(line, &tweet)
		if err != nil {
			log.Fatal("tweet parsed fail\n")
		}
		fmt.Println(string(line))
		fmt.Printf("[%s] %s\n", tweet.Id, tweet.Text)
	}
}
