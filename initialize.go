package linebot2lambda

import (
	"fmt"
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
	googleSecretJSON := os.Getenv("GOOGLE_SECRET_JSON")
	googleCalendarIDName := os.Getenv("GOOGLE_CALENDAR_ID")
	lineChannelAccessTokenName := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	lineChannelSecretName := os.Getenv("LINE_CHANNEL_SECRET")
	envMap, err := initParameter(
		googleSecretJSON,
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

	srv, err := initCalendarService([]byte(envMap[googleSecretJSON]))
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
	googleSecretJSON,
	googleCalendarIDName,
	lineChannelAccessTokenName,
	lineChannelSecretName string,
) (map[string]string, error) {
	svc := ssm.New(
		session.Must(session.NewSession()),
		aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")),
	)
	input := &ssm.GetParametersInput{
		Names: []*string{
			&googleSecretJSON,
			&googleCalendarIDName,
			&lineChannelAccessTokenName,
			&lineChannelSecretName,
		},
		WithDecryption: aws.Bool(true),
	}
	output, err := svc.GetParameters(input)
	if err != nil {
		err := fmt.Errorf("ssm GetParameters error occurred: %w", err)
		return nil, err
	}
	envMap := make(map[string]string)
	for _, p := range output.Parameters {
		envMap[*p.Name] = *p.Value
	}
	return envMap, nil
}

func initCalendarService(googleSecretJSON []byte) (*calendar.Service, error) {
	config, err := google.JWTConfigFromJSON(googleSecretJSON, calendar.CalendarScope)
	if err != nil {
		err := fmt.Errorf("Unable to parse client secret file to config: %w", err)
		return nil, err
	}
	client := config.Client(oauth2.NoContext)
	srv, err := calendar.New(client)
	if err != nil {
		err := fmt.Errorf("Unable to retrieve Calendar client: %w", err)
		return nil, err
	}
	return srv, nil
}
