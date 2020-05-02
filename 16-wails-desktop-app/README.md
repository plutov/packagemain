# packagemain #6: Building Desktop App in Go using Wails

Hi Gophers, My name is Alex Pliutau.

As we all know, Go is mostly used to build APIs, web backends, CLI tools. But what's interesting is that Go can be used in places we were not expecting to see it.

In this video we will build a Desktop applicaiton with Go and Vue.js using [Wails](https://wails.app) framework.

This framework is new and still in beta, but I was surprised how easy it was to develop, build and package an app with it.

Wails provides the ability to wrap both Go code and a web frontend into a single binary. The Wails CLI makes this easy for you, by handling project creation, compilation and bundling.

## App

I will build a very simple app to display CPU Usage of my machine in real time. And if you have time and like Wails, you can come up with something more creative and complex.

## Installation

Wails CLI can be installed with `go get`. After installation, you should set it up using `wails setup` command.

```
go get github.com/wailsapp/wails/cmd/wails
wails setup
```

Then let's bootstrap our project with the name `cpustats`:

```
wails init
cd cpustats
```

Our project consists of Go backend and Vue.js frontend. `main.go` will be our entrypoint, in which we can include any other dependencies, there is also `go.mod` file to manage them. `frontend` folder contains Vue.js components, webpack and CSS.

## Concepts

There are 2 main components to share data between Backend and Frontend: Binding and Events.

Binding is a single method that allows you to expose (bind) your Go code to the frontend. 

Also Wails provides a unified Events system similar to Javascript's native events system. This means that any event that is sent from either Go or Javascript can be picked up by either side. Data may be passed along with any event. This allows you to do neat things like have background processes running in Go and notifying the frontend of any updates.

## Backend

Let's develop a backend part first, to get CPU Usage and send it to the frontend using `bind` method.

I will create a new package and define a type which I'll expose (bind) to the frontend.

pkg/sys/sys.go:
```
package sys

import (
	"math"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/wailsapp/wails"
)

// Stats .
type Stats struct {
	log *wails.CustomLogger
}

// CPUUsage .
type CPUUsage struct {
	Average int `json:"avg"`
}

// WailsInit .
func (s *Stats) WailsInit(runtime *wails.Runtime) error {
	s.log = runtime.Log.New("Stats")
	return nil
}

// GetCPUUsage .
func (s *Stats) GetCPUUsage() *CPUUsage {
	percent, err := cpu.Percent(1*time.Second, false)
	if err != nil {
		s.log.Errorf("unable to get cpu stats: %s", err.Error())
		return nil
	}

	return &CPUUsage{
		Average: int(math.Round(percent[0])),
	}
}
```

If your struct has a `WailsInit` method, Wails will call it at startup. This allows you to do some initialisation before the main application is launched.

Import `sys` package in `main.go` and bind Stats instance to frontend:

```go
package main

import (
	"github.com/leaanthony/mewn"
	"github.com/plutov/packagemain/cpustats/pkg/sys"
	"github.com/wailsapp/wails"
)

func main() {
	js := mewn.String("./frontend/dist/app.js")
	css := mewn.String("./frontend/dist/app.css")

	stats := &sys.Stats{}

	app := wails.CreateApp(&wails.AppConfig{
		Width:  512,
		Height: 512,
		Title:  "CPU Usage",
		JS:     js,
		CSS:    css,
		Colour: "#131313",
	})
	app.Bind(stats)
	app.Run()
}
```

## Frontend

We bind the `stats` instance from Go, which can be used in frontend by callind `window.backend.Stats`. If we want to call a function `GetCPUUsage()` it will return us a Promise.

```js
window.backend.Stats.GetCPUUsage().then(cpu_usage => {
    console.log(cpu_usage);
})
```

To build the whole project into single binary we should run `wails build`, `-d` flag can be added to build a debuggable version. It will create a binary with a name matching the project name.

Let's test if it works by simply displaying the CPU Usage value on the screen:

```
wails build -d
./build/cpustats
```

## Events

We sent CPU Usage value to frontend using Binding, now let's try different approach, let's create a timer on Backend which will send CPU Usage values in the background using Events approach. Then we can subscribe to the event in Javascript.

In Go we can do it in `WailsInit` function:

```go
func (s *Stats) WailsInit(runtime *wails.Runtime) error {
	s.log = runtime.Log.New("Stats")

	go func() {
		for {
			runtime.Events.Emit("cpu_usage", s.GetCPUUsage())
			time.Sleep(1 * time.Second)
		}
	}()

	return nil
}
```

In Vue.js we can subscribe to this event when component is mounted (or any other place):

```js
mounted: function() {
    wails.Events.on("cpu_usage", cpu_usage => {
        if (cpu_usage) {
            console.log(cpu_usage.avg);
        }
    });
}
```

## Gauge Bar 

I'd like to display CPU Usage with a gauge bar, so I will include a third party dependency for that, simply by using `npm`:

```
cd frontend
npm install --save apexcharts
npm install --save vue-apexcharts
```

Then import it to `main.js` file:

```js
import VueApexCharts from 'vue-apexcharts'

Vue.use(VueApexCharts)
Vue.component('apexchart', VueApexCharts)
```

Now we can display our CPU Usage using apexcharts, and update the values of the component by receiving an event from Backend:

```html
<template>
  <apexchart type="radialBar" :options="options" :series="series"></apexchart>
</template>

<script>
export default {
  data() {
    return {
      series: [0],
      options: {
        labels: ['CPU Usage']
      }
    };
  },
  mounted: function() {
    wails.Events.on("cpu_usage", cpu_usage => {
      if (cpu_usage) {
        this.series = [ cpu_usage.avg ];
      }
    });
  }
};
</script>
```

To change styles we can directly modify the `src/assets/css/main.css` or define them in components.

## Final Build and Run

```
wails build -d
./build/cpustats
```

## Conclusion

I really enjoyed working with `Wails`, and the Events concept makes it really easy to control your application's state.

Check it out at [wails.app](https://wails.app) or on Github at [github.com/wailsapp/wails](https://github.com/wailsapp/wails)