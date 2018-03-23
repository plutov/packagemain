### Building Google Home Action with Go

Google Home is a voice Assistant, similar to Amazon Alexa, but working with Google services. It has a lot of built-in integrations, but what is interesting for us developes is that we can build our our programs for it. Google call them Actions.

In this video we will build an Action, which will help user to find an air quality index of the city user is located in. It's not necessary to have Google Home device to be able to build and test it, Google has very nice Similuator. However, I have device, so I'll show you how it works after we build it.

Google Home Actions are using DialogFlow (previously api.ai) to setup conversation flow using NLU. And we will build a simple backend API to get data, which we'll deploy to Google Cloud.

### Let's start.

We have to login with our Google Account to https://dialogflow.com. Go to Console and create your first project. You can choose to use existing Google Cloud project or create a new one.

We will use DialogFlow API V1, V2 is slightly different in terms of request/response format.

Now let's decide the future user flow.

There are 2 ways to start to talk to our Action: explicit invocation and implicit invocation. Explicit one is triggered when we tell Google "Talk to <action name>". In implicit invocation we can set up custom messages, but we will skip this option for our demo program. Basically you need to create an intent and describe possible sentences user may say.

We need to know user's location to get information, so first of all we need to ask for this permission. Google Action has functionality to ask for location permission. We need to send specific response to DialogFlow after user started to talk to Action.

Let's define it in welcome intent. We need to set the action name `location_permission` and then in our webhook we can check it. Also we need to `Enable webhook call for this intent`.

Let's describe our fallback intent with default fallback message. This intent will be executed when action doesn't understand what user wants.

`Sorry, I can't help you with this right now. Please try later.`. Set intent as end of conversation.

Now let's define the main intent to get air quality. This intent will be triggered not by specific word but by reserved event: `actions_intent_PERMISSION`. So when user granted access to location info this intent will be executed. We set action name as `get` and will handle it later in API. Also we need to enable webhook call for this. And set end of conversation.

### Ok, we're almost done with configuration, let's define how Google Action will get data.

Go to Fulfillment. There are 2 options to write a backend logic: using custom webhook or to use inline editor powered by cloud functions on Firebase, but it's node.js so we will go with first one and provide endpoint to API deployed to google cloud. We just need to set endpoint. I will do it later.

### Go

I already have Google Cloud project and `gcloud` SDK installed, so I will start with writing application. I will start with `app.yaml` file to describe handlers and runtime:

```go
runtime: go
api_version: go1

handlers:
- url: /
  script: _go_app
```

Let's create empty API and deploy it to see if it works:

```go
package app

import (
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/", handle)
}

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}
```

```
gcloud app deploy
```