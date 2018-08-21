package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/linebot"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

const (
	timeParseForm = "2006-01-02 15:04"
	timeZone      = "Asia/Tokyo"
)

var calendarId = os.Getenv("GOOGLE_CALENDAR_ID")

type Webhook struct {
	Events []Event `json:"events"`
}

type Event struct {
	ReplyToken string   `json:"replyToken"`
	Type       string   `json:"type"`
	Source     *Source  `json:"source"`
	Message    *Message `json:"message"`
}

type Source struct {
	Typev  string `json:"type"`
	UserId string `json:"userId"`
}

type Message struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	Text string `json:"text"`
}

func HandleRequest(events Webhook) (string, error) {
	log.Printf("%+v", events)
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)

	if err != nil {
		log.Fatal(err)
	}

	for _, event := range events.Events {
		log.Printf("%+v", event)
		if event.Type == "message" {
			log.Printf("%+v", event.Message.Text)
			if strings.HasPrefix(event.Message.Text, "予定登録\n") {
				b, err := ioutil.ReadFile("client_secret.json")
				if err != nil {
					log.Fatalf("Unable to read client secret file: %v", err)
				}

				config, err := google.JWTConfigFromJSON(b, calendar.CalendarScope)
				if err != nil {
					log.Fatalf("Unable to parse client secret file to config: %v", err)
				}
				client := config.Client(oauth2.NoContext)

				srv, err := calendar.New(client)
				if err != nil {
					log.Fatalf("Unable to retrieve Calendar client: %v", err)
				}
				calEvent := &calendar.Event{}
				msgs := strings.Split(event.Message.Text, "\n")
				log.Printf("%+v", msgs)
				for i, msg := range msgs {
					switch i {
					case 1:
						calEvent.Summary = msg
					case 2:
						calEvent.Location = msg
					case 3:
						calEvent.Description = msg
					case 4:
						start, _ := time.Parse(timeParseForm, msg)
						calEvent.Start = &calendar.EventDateTime{
							DateTime: start.Format(time.RFC3339),
							TimeZone: timeZone,
						}
					case 5:
						end, _ := time.Parse(timeParseForm, msg)
						calEvent.End = &calendar.EventDateTime{
							DateTime: end.Format(time.RFC3339),
							TimeZone: timeZone,
						}
					default:
						continue
					}
				}
				resultEvent, err := srv.Events.Insert(calendarId, calEvent).Do()
				if err != nil {
					log.Fatalf("%+v", err)
				}
				log.Printf("%+v", resultEvent)
				reply := "予定登録できたよ！\n" + resultEvent.HtmlLink
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
					log.Print(err)
				}
			}
		}

	}

	return fmt.Sprintf(""), nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(HandleRequest)
}
