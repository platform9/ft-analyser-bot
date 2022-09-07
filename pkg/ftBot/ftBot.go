package ftBot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	amplitudeapi "github.com/platform9/ft-analyser-bot/pkg/amplitude"
	"github.com/platform9/ft-analyser-bot/pkg/api"
	"github.com/platform9/ft-analyser-bot/pkg/config"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"go.uber.org/zap"
)

func FtBotRun() {

	//Load the slack creds
	slackCreds := config.SlackCreds()

	token := slackCreds.AuthToken
	appToken := slackCreds.AppToken

	// Create a new client to slack by giving token
	// Set debug to true while developing
	client := slack.New(token, slack.OptionDebug(true), slack.OptionAppLevelToken(appToken))
	socketClient := socketmode.New(
		client,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	// Create a context that can be used to cancel goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(ctx context.Context, client *slack.Client, socketClient *socketmode.Client) {
		// Create a for loop that selects either the context cancellation or the events incomming
		for {
			select {
			// inscase context cancel is called exit the goroutine
			case <-ctx.Done():
				log.Println("Shutting down socketmode listener")
				return
			case event := <-socketClient.Events:
				// More cases can be added here depending on the type of events this bot can respond to
				switch event.Type {
				// handle EventAPI events
				case socketmode.EventTypeEventsAPI:
					// The Event sent on the channel is not the same as the EventAPI events so we need to type cast it
					eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
					if !ok {
						log.Printf("Could not type cast the event to the EventsAPIEvent: %v\n", event)
						continue
					}
					socketClient.Ack(*event.Request)
					err := handleEventMessage(eventsAPIEvent, client)
					if err != nil {
						// Replace with actual err handeling
						log.Fatal(err)
					}
				}

			}
		}
	}(ctx, client, socketClient)

	socketClient.Run()
}

// handleEventMessage will take an event and handle it properly based on the type of event
func handleEventMessage(event slackevents.EventsAPIEvent, client *slack.Client) error {
	switch event.Type {
	// First we check if this is an CallbackEvent
	case slackevents.CallbackEvent:

		innerEvent := event.InnerEvent
		// Yet Another Type switch on the actual Data to see if its an AppMentionEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			// The application has been mentioned since this Event is a Mention event
			err := handleAppMentionEvent(ev, client)
			if err != nil {
				return err
			}
		case *slackevents.MessageEvent:
			// The application is a Message Event
			err := handleAppMessageEvent(ev, client)
			if err != nil {
				return err
			}
		}

	default:
		return errors.New("unsupported event type")
	}
	return nil
}

// Handles app mention (@bot) events for Weekly Analysis Events
func handleAppMentionEvent(event *slackevents.AppMentionEvent, client *slack.Client) error {
	// In case it is a null string don't respond
	if event.Text == "" && len(event.Text) == 0 {
		return nil
	}

	text := strings.ToLower(event.Text)

	// Create the attachment and assigned based on the message
	attachment := slack.Attachment{}

	dt := time.Now()

	// Weekly Analysis for FT Users
	weeklyAnalysis, err := amplitudeapi.WeeklyMessage()
	if err != nil {
		zap.S().Errorf("error while fetching weekly FT analysis, error: %v", err)
	}
	out := api.GenerateOutputString(weeklyAnalysis)

	// Add Some default context like user who mentioned the bot
	attachment.Fields = []slack.AttachmentField{
		{
			Title: "Date",
			Value: dt.Format("01-02-2006 15:04:05"),
		},
		{
			Title: "Weekly Analysis",
			Value: out,
		},
	}
	if strings.Contains(text, "weekly analysis") {
		attachment.Pretext = "FT Weekly Analysis"
		attachment.Color = "#FFA500"
	}
	// Send the message to the channel
	// The Channel is available in the event message
	_, _, err1 := client.PostMessage(event.Channel, slack.MsgOptionAttachments(attachment))
	if err1 != nil {
		return fmt.Errorf("failed to post message: %w", err1)
	}
	return nil
}

// Handles direct channel messages for NPS events
func handleAppMessageEvent(event *slackevents.MessageEvent, client *slack.Client) error {
	var userID string
	// In case it is a null string don't respond
	if event.Text == "" && len(event.Text) == 0 {
		return nil
	}

	// Check if the NPS score is greater than 8 and return from this function
	if strings.Contains(event.Text, ": `10`") || strings.Contains(event.Text, ": `9`") {
		return nil
	}

	text := strings.ToLower(event.Text)

	// hardcoding this for now
	res := event.Text

	// Just for our convinience, we are checking only for the NPS Apcues survey format and rejecting all other inputs
	if strings.Contains(res, "\n") && strings.Count(res, "\n") == 3 {
		res1 := strings.Split(res, "\n")
		userID = res1[1][14:]
	} else {
		return nil
	}

	// Create the attachment and assigned based on the message
	attachment := slack.Attachment{}

	dt := time.Now()
	npsAnalysis, err := amplitudeapi.NpsScoreAnalysis(userID, "")
	if err != nil {
		zap.S().Errorf("error while fetching nps analysis for user, error: %v", userID, err)
	}
	// Use this to prety print json
	/*jsonResp, err := json.MarshalIndent(npsAnalysis, "", "  ")
	if err != nil {
		fmt.Println(err)
		zap.S().Errorf("Error while marshalling the response: %v", err)
	}*/
	out := api.GenNPSOutput(npsAnalysis)

	// Add Some default context like user who mentioned the bot
	attachment.Fields = []slack.AttachmentField{
		{
			Title: "Date",
			Value: dt.Format("01-02-2006 15:04:05"),
		},
		{
			Title: "NPS Analysis",
			Value: string(out),
		},
	}

	if strings.Contains(text, "nps score") && strings.Contains(event.Text, "https://studio.appcues.com/nps|NPS Survey") {
		attachment.Pretext = "Analyzing the errors that user might have faced"
		attachment.Color = "#A020F0"
	}

	// Send the message to the channel
	// The Channel is available in the event message
	_, _, err1 := client.PostMessage(event.Channel, slack.MsgOptionAttachments(attachment))
	if err1 != nil {
		return fmt.Errorf("failed to post message: %w", err1)
	}

	return nil
}
