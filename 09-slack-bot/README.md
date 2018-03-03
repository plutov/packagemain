### Building a Slack Bot with Go and Wit.ai

In this video we will build a simple Slack Bot with NLU functionality to get some useful information from Wolfram. No worries if you didn't use Wolfram before, it's a computational knowledge engine which can give you a short answer to your question.

There are different platforms for NLU, such as LUIS.ai, Wit.ai, RASA_NLU, we will use Wit.ai, an NLU platform acquired by Facebook, provides you functionality to parse text and extract useful information.

Our Bot will have minimal functionality, such as: "Reply to user on greeting" and then "Search in Wolfram and reply back". Why do we need NLU here? Because we don't know how user greets us, also we don't know how user will ask our Bot, will it be something like "Do you know who is the president of USA?" or "Can you tell me who is the president of USA.". Also it could be a typo from user. And NLU will help us to extract useful info from a custom message.

Let's start with creating our new Bot in Slack, for this I created new Slack team.

 - https://packagemain.slack.com/apps
 - Search for Bots
 - Add Configuration
 - Username "packagemain"
 - Configure other settigns if you need and get token
 - Now we can check if bot exists: Hi

Let's write some simple program which will use Slack real time messaging to receive user input. We need to use Slack token, which I already set locally.

Let's handle each incoming message in separate go routine.

```
go get github.com/nlopes/slack
```

```
var (
	slackClient   *slack.Client
)

func main() {
	slackClient = slack.New(os.Getenv("SLACK_ACCESS_TOKEN"))

	rtm := slackClient.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			go handleMessage(ev)
		}
	}
}

func handleMessage(ev *slack.MessageEvent) {
	fmt.Printf("%v\n", ev)
}
```

#### Setup Wit.ai

Before we're going to make Bot reply to user, we should configure Wit.ai. We're going to create an application there now.

 - Go to https://wit.ai/home
 - New App
 - Define Entities

Wit.ai has predefined entities and we will use 2 of them. We can also define our own and train Wit.ai to understand it, but I'll leave it to you as a homework.

1. wit/greetings: Hi, Hello.
2. wit/wolfram_search_query: Who is the president of Belarus, distance between Earth and Mars, formula of ethanol.

I searched for Wit.ai Go package, there are 5 packages on godoc.org and 4 of them are not compatible with new API. And one which is working has 2 stars on GitHub. We will use it, but I don't recommend to use it in Production environments. Wanna try to create Go package - create Wit.ai SDK please.

Get server access token in Settings

```
go get github.com/christianrondeau/go-wit
```

```
result, err := witClient.Message(ev.Msg.Text)
if err != nil {
	log.Printf("unable to get wit.ai result: %v", err)
	return
}

fmt.Printf("%v\n", result)
```

Wit.ai will return us list of entities, each entity contains value and confidence, so let's filter entities with low confidence, for that we will set a confidence threshold as 0.5.

```
const confidenceThreshold = 0.5
```

As you can see when we type "hi" Wit.ai returns us 2 successfull entities as "hi" is a valid wolfram search also. But greetings entity has a higher confidence, let's use one with highest confidence.

```
var (
	topEntity    wit.MessageEntity
	topEntityKey string
)

for key, entityList := range result.Entities {
	for _, entity := range entityList {
		if entity.Confidence > confidenceThreshold && entity.Confidence > topEntity.Confidence {
			topEntity = entity
			topEntityKey = key
		}
	}
}
```

Now we can reply to user based on input:

```
func replyToUser(ev *slack.MessageEvent, topEntity wit.MessageEntity, topEntityKey string) {
	switch topEntityKey {
	case "greetings":
		slackClient.PostMessage(ev.User, "Hello user! How can I help you?", slack.PostMessageParameters{
			AsUser: true,
		})
		return
	case "wolfram_search_query":
		// TODO
	}

	slackClient.PostMessage(ev.User, "¯\\_(o_o)_/¯", slack.PostMessageParameters{
		AsUser: true,
	})
}
```

We can see that our program also handles messages send by itself, let's fix that.

```
if len(ev.BotID) == 0 {
	go handleMessage(ev)
}
```

#### Setup Wolfram

We can say hello, let's now get answer to user's question from Wolfram:

 - Go to https://developer.wolframalpha.com/portal/myapps
 - Click Get App ID
 - Copy APP ID

```
go get github.com/Krognol/go-wolfram
```

```
res, err := wolframClient.GetSpokentAnswerQuery(topEntity.Value.(string), wolfram.Metric, 1000)
if err == nil {
	slackClient.PostMessage(ev.User, res, slack.PostMessageParameters{
		AsUser: true,
	})
	return
}

log.Printf("unable to get data from wolfram: %v", err)
```

#### Testing part

Hi
Hola

Who is the president of US?

What is the meaning of life?

bye