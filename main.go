package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/linebot"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

const (
	timeParseForm = "2006-01-02 15:04"
	timeZone      = "Asia/Tokyo"
)

var calendarId = os.Getenv("GOOGLE_CALENDAR_ID")

var (
	location        *time.Location
	calendarService *calendar.Service
	linebotClient   *linebot.Client
)

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
	Type    string `json:"type"`
	UserId  string `json:"userId"`
	GroupId string `json:"groupId"`
	RoomId  string `json:"roomId"`
}

type Message struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	Text string `json:"text"`
}

func HandleRequest(events Webhook) error {
	for _, event := range events.Events {
		fmt.Printf("%#v\n", event)
		if event.Type != "message" {
			continue
		}

		if strings.HasPrefix(event.Message.Text, "予定登録\n") {
			registerSchedule(event)
		} else if strings.HasPrefix(event.Message.Text, "予定登録フォーマット") {
			replyTemplate(event)
		}
	}
	return nil
}

func registerSchedule(event Event) {
	msgs := strings.Split(event.Message.Text, "\n")
	log.Printf("%+v", msgs)
	start, _ := time.ParseInLocation(timeParseForm, msgs[4], location)
	end, _ := time.ParseInLocation(timeParseForm, msgs[5], location)
	calEvent := &calendar.Event{
		Summary:     msgs[1],
		Location:    msgs[2],
		Description: msgs[3],
		Start: &calendar.EventDateTime{
			DateTime: start.Format(time.RFC3339),
			TimeZone: timeZone,
		},
		End: &calendar.EventDateTime{
			DateTime: end.Format(time.RFC3339),
			TimeZone: timeZone,
		},
	}
	resultEvent, err := calendarService.Events.Insert(calendarId, calEvent).Do()
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
	reply := "予定登録できたよ！\n" + resultEvent.HtmlLink
	if _, err = linebotClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
		log.Printf("%v\n", err)
		return
	}
}

func replyTemplate(event Event) {
	reply := "予定登録フォーマット↓だよ。１行目は必ず`予定登録`っていれてね。\n\n予定登録\n[タイトル]\n[場所]\n[詳細]\n[開始時間(2018-01-02 12:30)]\n[終了時間(2018-01-03 20:30)]"
	if _, err := linebotClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
		log.Printf("%v\n", err)
		return
	}
}

func initLocation() {
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	location = loc
}

func initCalendarService() {
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
	calendarService = srv
}

func initLinebotClient() {
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	linebotClient = bot
}

func main() {
	initLocation()
	initCalendarService()
	initLinebotClient()
	lambda.Start(HandleRequest)
}
