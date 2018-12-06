## packagemain #15: Dialogflow v2 SDK for Go

Hi Gophers, My name is Alex Pliutau. Welcome to package main, the channel about Go.

As you may remember I already have a [video](https://www.youtube.com/watch?v=LeGuJo7QBbI) about building Google Home Action with Go. But at that time DialogFlow had no Go SDK so we had to make manuall HTTP requests and define our DialogFlow types. Now Go is listed on [official page](https://dialogflow.com/docs/sdks) of V2 clients. However if we go to [GoDoc](https://godoc.org/cloud.google.com/go/dialogflow/apiv2) of this SDK it says "NOTE: This package is in alpha. It is not stable, and is likely to change.". Not a problem, we will still give it a try.

One subscriber asked me to show how to make a Google Action to grab bitcoin price from [coinmarketcap.com](https://coinmarketcap.com/), and I thought it's a good time to try this new SDK.

### And this is how we're going to make it

First of all we need to define intents in Dialogflow, in our case user will say "Hey Google. Please tell me the price of Bitcoin.". User may ask for the price of Litecoin or other crypto currency, so we have to get it from the input, Dialogflow will help us here. When we know the name of the currency, we can call [coinmarketcap API](https://coinmarketcap.com/api/) to get the price. Of course we will handle errors and edge cases if currency is not found, or there is some issue with API call. And if all good we will ask Dialogflow to respond user with the price.

Sounds good? Let's make it!

### Public endpoint

For Dialogflow we should have external endpoint of our service, so I am going to deploy our application to Google App Engine.

#### app.yaml

```yaml
runtime: go
api_version: go1

handlers:
- url: /
  script: _go_app
```