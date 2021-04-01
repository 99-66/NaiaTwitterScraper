package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	METHOD = "GET"
	URL = "https://api.twitter.com/2/tweets/sample/stream"
	TOKEN = ""
)

type Tweet struct {
	Id string `json:"id"`
	Text string `json:"text"`
}

type StreamTweet struct {
	Tweet `json:"data"`
}

func main() {
	client := &http.Client{}
	req, err := http.NewRequest(METHOD, URL, nil)
	if err != nil {
		log.Fatalf("http new request failed. %v\n", err)
	}

	bearerToken := "Bearer " + TOKEN
	req.Header.Add("Authorization", bearerToken)

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