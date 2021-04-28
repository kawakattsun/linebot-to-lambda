package linebot2lambda

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/line/line-bot-sdk-go/linebot"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

const timeZone = "Asia/Tokyo"

type Config struct {
	Location         *time.Location
	CalendarService  *calendar.Service
	LinebotClient    *linebot.Client
	GoogleCalendarID string
}

func Initialize() (*Config, error) {
	googleCalendarIDName := os.Getenv("GOOGLE_CALENDAR_ID")
	lineChannelAccessTokenName := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	lineChannelSecretName := os.Getenv("LINE_CHANNEL_SECRET")
	envMap, err := initParameter(
		googleCalendarIDName,
		lineChannelAccessTokenName,
		lineChannelSecretName,
	)
	if err != nil {
		return nil, err
	}
	c := new(Config)
	c.GoogleCalendarID = envMap[googleCalendarIDName]

	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return nil, err
	}
	c.Location = loc

	srv, err := initCalendarService()
	if err != nil {
		return nil, err
	}
	c.CalendarService = srv

	bot, err := linebot.New(
		envMap[lineChannelSecretName],
		envMap[lineChannelAccessTokenName],
	)
	if err != nil {
		return nil, err
	}
	c.LinebotClient = bot

	return c, nil
}

func initParameter(
	googleCalendarIDName,
	lineChannelAccessTokenName,
	lineChannelSecretName string,
) (map[string]string, error) {
	svc := ssm.New(
		session.Must(session.NewSession()),
		aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")),
	)
	input := ssm.GetParametersInput{
		Names: []*string{
			aws.String(googleCalendarIDName),
			aws.String(lineChannelAccessTokenName),
			aws.String(lineChannelSecretName),
		},
		WithDecryption: aws.Bool(true),
	}
	output, err := svc.GetParameters(input)
	if err != nil {
		fmt.Errorf("ssm GetParameters error occurred: %w", err)
		return nil, err
	}
	envMap := make(map[string]string)
	for _, p := range output.GetParameters {
		envMap[*p.Name] = *p.Value
	}
	return envMap, nil
}

func initCalendarService() (*calendar.Service, error) {
	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		fmt.Errorf("Unable to read client secret file: %w", err)
		return nil, err
	}
	config, err := google.JWTConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		fmt.Errorf("Unable to parse client secret file to config: %w", err)
		return nil, err
	}
	client := config.Client(oauth2.NoContext)
	srv, err := calendar.New(client)
	if err != nil {
		fmt.Errorf("Unable to retrieve Calendar client: %w", err)
	}
	return srv, nil
}
