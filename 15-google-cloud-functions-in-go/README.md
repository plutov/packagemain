# packagemain #15: Google Cloud Functions in Go

Hi Gophers, My name is Alex Pliutau.

Earlier this month Google Cloud Functions team finally announced beta support of Go, the runtime uses Go 1.11, which includes go modules as we know.

And in this video I am going to show how to write and deploy 2 types of functions: HTTP function and background function.

## 2 types of functions

HTTP functions are functions that are invoked by HTTP requests. They follow the http.HandlerFunc type from the standard library.

In contrast, background functions are triggered in response to an event. Your function might, for example, run every time there is new message in Cloud Pub/Sub service.

## Setup Google Cloud Project

The first step is to ensure that you have a Google Cloud Platform account with Billing setup. Remember that there is an always free tier when you sign up for Google Cloud and you could use that too.

Once you have setup your project, the next step is to enable the Google Cloud Functions API for your project. You can do it from Cloud Console or from your terminal using `gcloud` tool.

```bash
gcloud services enable cloudfunctions.googleapis.com
```

## HTTP function

It will be a simple HTTP function which generates a random number and sends it to Cloud Pub/Sub topic.

Let's create our topic first:

```bash
gcloud pubsub topics create randomNumbers
```

I will create a separate folder / package for this function.

CODE 1

Our package uses `cloud.google.com/go/pubsub` package, so let's initialize go modules.

```bash
export GO111MODULE=on
go mod init
go mod tidy
```

If you have external dependencies, you have to vendor them under the library package locally before deploying.

```bash
go mod vendor
```

Now it's time to deploy it:

```bash
gcloud functions deploy api --entry-point Send --runtime go111 --trigger-http --set-env-vars PROJECT_ID=projectname-227718
```

Where `api` is a name, `Send` is an entrypoint function, `--trigger-http` tells that it is HTTP function. And we also set a PROJECT_ID env var.

The deployment may take few minutes.

HTTP functions can be reached without an additional API gateway layer. Cloud Functions give you an HTTPS URL. After the function is deployed, you can invoke the function by entering the URL into your browser.

Output:

```bash
availableMemoryMb: 256
entryPoint: Send
environmentVariables:
  PROJECT_ID: projectname-227718
httpsTrigger:
  url: https://us-central1-projectname-227718.cloudfunctions.net/api
```

## Background function

Since Google Cloud background functions can be triggered from Pub/Sub, let's just write a function which will simply log a payload of event triggering it.

CODE 1

Note: we don't need go modules in consumer.

The deployment part is very similar to HTTP function, except how we're triggering this function.

```bash
gcloud functions deploy consumer --entry-point Receive --runtime go111 --trigger-topic=randomNumbers
```

Let's check logs now after execution.

```bash
gcloud functions logs read consumer
```

## Cleanup

To cleanup let's delete everything we created: function and pub/sub topic.

```bash
gcloud functions delete api
gcloud functions delete consumer
gcloud pubsub topics delete randomNumbers
```

## Conclusion

Please share your experience with Google Cloud Functions in Go, are you missing any functionality there, any issues you encountered.