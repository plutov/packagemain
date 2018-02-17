### Code review of push notification server

Repository for review: https://github.com/ortuman/mercury - Push notification server written in Go.

```
go get github.com/ortuman/mercury
git co -b codereview
```

To make a valid test of our changes we need to have this server running locally.

Let's run it with go run:
```
go run mercury.go --config example.mercury.conf
```

#### mercury.go:

 - There is no exit code -1: `os.Exit(1)`
 - Also we can remove `else`
 - We can remove identation by checking `help` first
 - In first `if` we use fmt+Exit, but later we use log.Fatalf. I think we can use everywhere log.Fatal

Test:
```
go run mercury.go
go run mercury.go --help
```

#### config/config.go

First thing I don't like is that we have a lot of global variables. Instead of assigning config to global vars I think it's better if `config.Load` function will return us config object and then we can pass it to the server.

 - `Load` returns `globalConfig` -> `Config`
 - `initDefaultSettings` -> `getDefaultConf`
 - Remove `init`
 - Remove global vars

```
// Load config file
func Load(cfgFile string) Config {
	var conf Config
	if _, err := toml.DecodeFile(cfgFile, &conf); err != nil {
		logger.Warnf("config: couldn't load config file '%s': %v", cfgFile, err)
		return getDefaultConf()
	}

	return conf
}

func getDefaultConf() globalConfig {
	return globalConfig{
		Logger: LoggerConfig{
			Level:   "DEBUG",
			Logfile: "mercury.log",
		},
		Server: ServerConfig{
			ListenAddr: ":8080",
		},
		Redis: RedisConfig{
			Host: "localhost:6379",
		},
		Apns: ApnsConfig{
			MaxConn:         16,
			CertFile:        "cert.p12",
			SandboxCertFile: "cert.p12",
		},
		Gcm: GcmConfig{
			MaxConn: 16,
		},
		Safari: SafariConfig{
			MaxConn: 16,
		},
		Chrome: WebPushConfig{
			MaxConn: 16,
		},
		Firefox: WebPushConfig{
			MaxConn: 16,
		},
	}
}
```

#### mercury.go

```
conf := config.Load(configFile)

srv := server.NewServer(conf)
```

#### server/server.go

```
type Server struct {
	Config config.Config
}

func NewServer(conf config.Config) *Server {
	return &Server{
		Config: conf,
	}
}
```

Now a lot of places will break because we removed global vars, so we need to replace it everywhere.
Replace: `config.` -> `s.Config.`

 - Remove `else` from the end of `Run()` func.
 - Simplify `if-else` in `createPIDFile`