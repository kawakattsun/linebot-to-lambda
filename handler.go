package linebot2lambda

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"google.golang.org/api/calendar/v3"
)

func HandleRequest(c *Config, r *http.Request) error {
	events, err := c.LinebotClient.ParseRequest(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse request error occurred: %v\n", err)
		return nil
	}

	for _, event := range events {
		if event.Type != linebot.EventTypeMessage {
			continue
		}
		switch message := event.Message.(type) {
		case *linebot.TextMessage:
			if strings.HasPrefix(message.Text, "予定登録\n") {
				registerSchedule(c, message, event.ReplyToken)
			} else if strings.HasPrefix(message.Text, "予定登録フォーマット") {
				replyTemplate(c, event.ReplyToken)
			}
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

func registerSchedule(c *Config, message *linebot.TextMessage, replyToken string) {
	msgs := strings.Split(message.Text, "\n")
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
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	reply := "予定登録できたよ！\n" + resultEvent.HtmlLink
	if _, err = c.LinebotClient.ReplyMessage(replyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
}

func replyTemplate(c *Config, replyToken string) {
	reply := "予定登録フォーマット↓だよ。１行目は必ず`予定登録`っていれてね。\n\n予定登録\n[タイトル]\n[場所]\n[詳細]\n[開始時間(2018-01-02 12:30)]\n[終了時間(2018-01-03 20:30)]"
	if _, err := c.LinebotClient.ReplyMessage(replyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
}
