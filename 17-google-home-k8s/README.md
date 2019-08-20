## packagemain #17: Google Home Action to manage your Kubernetes cluster

I always wanted to find a good use case of Google Home to make some DevOps tasks funnier. For example voice deployments, system metrics, etc. Since I use Kubernetes a lot, I thought it would be fun to control it via voice commands.

So I decided to develop this Action with the following goals in mind:

- Action should be able to work with any Kubernetes cluster: GKE, EKS, etc.
- Extensible commands
- Intuitive voice UX
- Easy installation

I published this project on Github ([google-home-k8s](https://github.com/plutov/google-home-k8s)), and I'll do a demo of it in this video.

### How it all works

The diagram below shows the flow from the user device to Dialogflow to our API deployed to App Engine (could be deployed somewhere else) to our Kubernetes cluster.

![diagram](https://raw.githubusercontent.com/plutov/packagemain/master/17-google-home-k8s/diagram.png)

Let's review some of the components:

- Dialogflow is a development suite for creating conversational interfaces. Here we will set up the user flow and specify all available voice commands, intents, and entities.
- Voice commands will call our API deployed to App Engine.
- Depends on the command, API will call Kubernetes API.

### Installation

#### kubeconfig

I assume that you already have a Kubernetes cluster running somewhere. I will use my demo GKE cluster in this video.

First of all, we have to generate a valid `kubeconfig` configuration file which will be used by our API for authentication.

Let's clone our project and use a script I prepared to generate the `kubeconfig`. Before doing this make sure that your `kubectl` points to right config and you have necessary permissions in this cluster.

```
git clone git@github.com:plutov/google-home-k8s.git
cd google-home-k8s
./generate-kubeconfig.sh
```

It will generate a `./build/kubeconfig` file which will be used by API later, this configuration never expires.

#### Deploy API

Now you can deploy the API somewhere, I will use Google App Engine for this purpose. API requires few env. variables to be set such as `NAMESPACE`, `API_KEY` and `LOG_LEVEL`.

```
cp env.sample.yaml env.yaml
```

API is protected by static API Key which should be set in `env.yaml`. To access API, the client should send `Authorization: Bearer ${API_KEY}` header.

```
gcloud app deploy
```

Copy the URL of your API.

#### Dialogflow

**google-home-k8s** already contains a Dialogflow configuration which can be easily imported.

1. Go to [Dialogflow Console](https://console.dialogflow.com/)
2. Select or create a new agent
3. Go to Settings -> Export and Import
4. Select **Import From Zip** (import this file [google-home-k8s.zip](https://raw.githubusercontent.com/plutov/google-home-k8s/master/google-home-k8s.zip))
5. Go to Fulfillment
6. Enable Webhook
7. Paste URL to API deployed to App Engine
8. Add Header. Key: `Authorization`, Value: `Bearer API_KEY` (replace `API_KEY` with the value from `env.yaml`)

#### Test

Go to Google Assistant, give your Action a name "Kubernetes Manager" and click "Test" which will open a Simulator.

If everything was done correctly you'll be able to reproduce the following conversation:

Example conversation:

> [you] Hey Google, talk to Kubernetes Manager

> [assistant] Welcome to Kubernetes Manager. How can I help you?

> [you] Scale statefulset "redis"

> [assistant] Got it. Currently, there are 3 replicas of the "redis" statefulset. To how many replicas do you want to scale?

> [you] 5

> [assistant] Statefulset has been updated. Anything else?

If you have Google Home device you can do same conversation there.

### Conclusion

I don't suggest to publish this Google Home Action, but rather use it from the single device authenticated with the same user as in Google Assistant. Unless you want to make your Kubernetes cluster public :)

Also, feel free to contribute to the project on Github to add more commands.