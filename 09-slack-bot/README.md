### Building a Slack Bot with Go and Wit.ai

In this video we will build a simple Slack Bot with NLU functionality to get some useful information from Wolfram.

Wit.ai is an NLU platform acquired by Facebook, provides you functionality to parse text and extract useful information.

#### Setup Wit.ai

 - New App
 - Define Entities
 - Get server access token

Wit.ai has predefined entities and we will use few of them.

1. wit/greetings: Hi, Hello.
2. wit/wolfram_search_query: president of Belarus, distance between Earth and Mars, formula of ethanol.

#### Setup Wolfram

 - Create new app
 - Get App ID

#### Setup Slack

 - Create new integration
 - Get access token

#### Development

We will use some external packages, so let's start with adding them.

```
dep init
dep ensure --add github.com/nlopes/slack
```

Go packages for Wit.ai and Wolfram are not so active, use them on your own risk:

```
dep ensure --add github.com/jsgoecke/go-wit
dep ensure --add github.com/Krognol/go-wolfram
```

I prepared all env variables on my host already.