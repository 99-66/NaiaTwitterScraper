package kafka

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/99-66/NaiaTwitterScraper/config"
	"github.com/99-66/NaiaTwitterScraper/models"
	"github.com/Shopify/sarama"
	"github.com/dghubble/oauth1"
	"io"
	"log"
	"os"
	"time"
)

func NewAsyncProducer(broker string) (sarama.AsyncProducer, error) {
	saramaCfg := sarama.NewConfig()
	producer, err := sarama.NewAsyncProducer([]string{broker}, saramaCfg)
	if err != nil {
		return nil, err
	}

	return producer, nil
}


// GetStreamTweets 트위터 샘플 스트림에서 트윗을 받아온다
// 트윗은 한글만 받는다
// 받아온 트윗은 채널로 전달된다
func GetStreamTweets(api *config.API, c chan<- models.Tweet) {
	URL := "https://stream.twitter.com/1.1/statuses/sample.json?language=ko"
	conf := oauth1.NewConfig(api.ApiKey, api.ApiSecret)
	token := oauth1.NewToken(api.AccessToken, api.AccessTokenSecret)

	httpClient := conf.Client(oauth1.NoContext, token)
	resp, err := httpClient.Get(URL)
	if err != nil {
		log.Fatalf("http request failed. %v\n", err)
	}
	defer resp.Body.Close()
	reader := bufio.NewReader(resp.Body)

	for {
		line, err := reader.ReadBytes('\r')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("buffer read bytes failed. %v\n", err)
			// 스트림 에러가 발생하면 딜레이를 준 후 종료시킨다
			// 종료되더라도 docker --restart 옵션으로 다시 실행된다
			time.Sleep(time.Second * 30)
			os.Exit(1)
		}

		var tweet *models.Tweet
		err = json.Unmarshal(line, &tweet)
		if err != nil {
			log.Printf("tweet parsing failed. %v\n", err)
			continue
		}
		tweet.TrimText()
		err = tweet.SetTag()
		if err != nil {
			log.Printf("set tagged failed. %v\n", err)
			continue
		}
		err = tweet.ChangeDateFormat()
		if err != nil {
			log.Printf("parsing to tweet failed. %v\n", err)
			continue
		}
		err = tweet.SetOrigin()
		if err != nil {
			log.Printf("set origin to failed. %v\n", err)
			continue
		}
		// 리트윗이라면 처리하지 않고 넘어간다
		if tweet.IsRetweet() {
			continue
		}
		c <- *tweet
	}
	defer close(c)
}

func TestGetStreamTweets(api *config.API, c chan<- models.Tweet) {
	for i:=0; i < 10; i ++{
		fmt.Println("goroutine ", i)
		tw := models.Tweet{
			CreatedAt: "2021-06-01T11:27:08+09:00",
			Id: 1377447408863342596,
			Text: "편한건가 싶어서 조금 후회된다 티셔츠 입을 걸",
		}
		err := tw.SetOrigin()
		if err != nil {
			log.Fatalf("set origin to failed. %v\n", err)
		}
		log.Printf("send to %v\n", tw)
		c <- tw
	}

	close(c)
}
