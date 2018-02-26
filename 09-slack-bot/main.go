package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jsgoecke/go-wit"
	"github.com/nlopes/slack"
)

func main() {
	witClient := wit.NewClient(os.Getenv("WIT_AI_ACCESS_TOKEN"))
	slackClient := slack.New(os.Getenv("SLACK_ACCESS_TOKEN"))

	rtm := slackClient.NewRTM()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			fmt.Printf("%v\n", ev)

			result, err := witClient.Message(&wit.MessageRequest{
				Query: ev.Msg.Text,
			})
			if err != nil {
				log.Printf("unable to get wit.ai result: %v", err)
				continue
			}

			fmt.Printf("%v\n", result)
		}
	}
}
