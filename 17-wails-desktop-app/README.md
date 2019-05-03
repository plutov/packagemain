# packagemain #17: Building Desktop App using Wails

Hi Gophers, My name is Alex Pliutau.

As we all know, Go is mostly used to build APIs, web backends, CLI tools. But what's interesting is that Go can be used in places we were not expecting to see it.

In this video we will build a Desktop applicaiton with Go and Vue.js using [wails.app](https://wails.app) framework.

This framework is new and still in beta, but I was surprised how easy it was to develop, build and package an app with it.

Wails provides the ability to wrap both Go code and a web frontend into a single binary. The Wails CLI makes this easy for you, by handling project creation, compilation and bundling.

## App

I will build a very simple app to display CPU Usage of my machine in real time. And if you have a time and like Wails, you can come up with something more creative and complex.

## Installation

Let's install Wails CLI using `go get` and set it up:

```
go get -u github.com/wailsapp/wails/cmd/wails
wails setup
```

Then let's bootstrap our project with the name `cpustats`:

```
wails init
cd cpustats
```

Our project consists of Go backend and Vue.js frontend. `main.go` will be our entrypoint, in which we can include any other dependencies. `frontend` folder contains Vue.js components, webpack and CSS.

## Concepts

There are 2 main components to share data between Backend and Frontend: Binding and Events.

Binding is a single method that allows you to expose (bind) your Go code to the frontend. 

Also Wails provides a unified Events system similar to Javascript's native events system. This means that any event that is sent from either Go or Javascript can be picked up by either side. Data may be passed along with any event. This allows you to do neat things like have background processes running in Go and notifying the frontend of any updates.

## Backend

Let's develop a backend part first, to get CPU Usage and send it to the frontend using binding or events.

## Frontend

I'd like to display CPU Usage with a gauge bar, so I wil include a third party dependency for that, simplpy using `npm`:

```
npm install --save apexcharts
npm install --save vue-apexcharts
```

To change styles we can directly modify the `src/assets/css/main.css` or define them in components.

## Build and Run

To build the whole project into single binary we should run `wails build`, `-d` flag can be added to build a debuggable version.

It will create a binary with a name matching the project name.

```
wails build -d
./cpustats
```
