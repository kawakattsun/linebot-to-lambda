package linebot2lambda

import (
	"log"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"google.golang.org/api/calendar/v3"
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

func HandleRequest(c *Config, events Webhook) error {
	for _, event := range events.Events {
		if event.Type != "message" {
			continue
		}

		if strings.HasPrefix(event.Message.Text, "予定登録\n") {
			registerSchedule(c, event)
		} else if strings.HasPrefix(event.Message.Text, "予定登録フォーマット") {
			replyTemplate(c, event)
		}
	}
	return nil
}

const timeParseForm = "2006-01-02 15:04"

const (
	calendarTitle = iota
	calendarSummary
	calendarLocation
	calendarDescription
	calendarStart
	calendarEnd
)

func registerSchedule(c *Config, event Event) {
	msgs := strings.Split(event.Message.Text, "\n")
	start, _ := time.ParseInLocation(timeParseForm, msgs[calendarStart], c.Location)
	end, _ := time.ParseInLocation(timeParseForm, msgs[calendarEnd], c.Location)
	calEvent := &calendar.Event{
		Summary:     msgs[calendarSummary],
		Location:    msgs[calendarLocation],
		Description: msgs[calendarDescription],
		Start: &calendar.EventDateTime{
			DateTime: start.Format(time.RFC3339),
			TimeZone: c.Location.String(),
		},
		End: &calendar.EventDateTime{
			DateTime: end.Format(time.RFC3339),
			TimeZone: c.Location.String(),
		},
	}
	resultEvent, err := c.CalendarService.Events.Insert(c.GoogleCalendarID, calEvent).Do()
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
	reply := "予定登録できたよ！\n" + resultEvent.HtmlLink
	if _, err = c.LinebotClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
		log.Printf("%v\n", err)
		return
	}
}

func replyTemplate(c *Config, event Event) {
	reply := "予定登録フォーマット↓だよ。１行目は必ず`予定登録`っていれてね。\n\n予定登録\n[タイトル]\n[場所]\n[詳細]\n[開始時間(2018-01-02 12:30)]\n[終了時間(2018-01-03 20:30)]"
	if _, err := c.LinebotClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
		log.Printf("%v\n", err)
		return
	}
}
