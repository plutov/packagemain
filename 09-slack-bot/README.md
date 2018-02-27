### Building a Slack Bot with Go and Wit.ai

In this video we will build a simple Slack Bot with NLU functionality to get some useful information from Wolfram. No worries if you didn't use Wolfram before, it's a computational knowledge engine which can give you a short answer to your question.

There are different platforms for NLU, such as LUIS.ai, Wit.ai, RASA_NLU, we will use Wit.ai, an NLU platform acquired by Facebook, provides you functionality to parse text and extract useful information.

Our Bot will have minimal functionality, such as: "Reply to user on greeting" and then "Search in Wolfram and reply back". Why do we need NLU here? Because we don't know how user greets us, also we don't know how user will ask our Bot, will it be something like "Do you know who is the president of USA?" or "Can you tell me who is the president of USA.". Also it could be a typo from user. And NLU will help us to extract useful info from a custom message.

Let's start with Wit.ai.

#### Setup Wit.ai

 - New App
 - Define Entities
 - Get server access token

Wit.ai has predefined entities and we will use 2 of them. We can also define our own and train Wit.ai to understand it, but I'll leave it to you as a homework.

1. wit/greetings: Hi, Hello.
2. wit/wolfram_search_query: Who is the president of Belarus, distance between Earth and Mars, formula of ethanol.

#### Setup Wolfram

 - Go to https://developer.wolframalpha.com/portal
 - Create new app
 - Get App ID

#### Setup Slack

 - Create new Custom integration
 - Get access token

#### Development

Let's start with Slack Real time messaging using Go slack package

```
go get github.com/nlopes/slack
```

I store all keys in `keys` file. We need 3 keys for our program: WOLFRAM_APP_ID, WIT_AI_ACCESS_TOKEN, SLACK_ACCESS_TOKEN.

Let's handle each message in a separate Go routine.

Then I searched for Wit.ai Go package, there are 5 packages on godoc.org and 4 of them are not compatible with new API. And one which is working has 2 stars on GitHub. We will use it, but I don't recommend to use it in Production environments. Wanna try to create Go package - create Wit.ai SDK please.

```
go get github.com/christianrondeau/go-wit
```

Wit.ai will return us list of entities, each entity contains value and confidence, so let's filter entities with low confidence, we will set a confidence threshold as 0.5.

When we type "hi" Wit.ai returns us 2 successfull entities as "hi" is a valid wolfram search also. But greetings entity has a higher confidence, let's use one with highest confidence.

Now we can see that our program also handles messages send by itself, let's fix that.

#### Wolfram

Let's get wolfram Go package.

```
go get github.com/Krognol/go-wolfram
```

#### Testing part

Hi
Hola

Who is the president of US?

What is the meaning of life?

bye