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
			log.Fatalf("buffer read bytes failed. %v\n", err)
		}

		var tweet *models.Tweet
		err = json.Unmarshal(line, &tweet)
		if err != nil {
			log.Fatalf("tweet parsing failed. %v\n", err)
		}
		tweet.TrimText()
		err = tweet.ChangeDateFormat()
		if err != nil {
			log.Fatalf("parsing to tweet failed. %v\n", err)
		}
		err = tweet.SetOrigin()
		if err != nil {
			log.Fatalf("set origin to failed. %v\n", err)
		}
		// {2021-04-01T11:27:08+09:00 1377447408863342594 납골당은 첨 가본돠 옷차림이나 뭐나 뚜렷한 형식은 없다길래 편하게 입어봤돠 넘 편한 건가 싶어서 대한유교걸 조금 후회된돠 슬랙스 입을 걸 그랬나??}
		log.Printf("send to %d %s\n", tweet.Id, tweet.Text)
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
			Text: "첨 가본돠 옷차림이나 뭐나 뚜렷한 형식은 없다길래 편하게 입어봤돠 넘 편한 건가 싶어서 대한유교걸 조금 후회된돠 슬랙스 입을 걸 그랬나??",
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
