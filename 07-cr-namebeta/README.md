### Code review of command line domain query tool

Repository: https://github.com/TimothyYe/namebeta
```
go get github.com/TimothyYe/namebeta
```

#### Common package review

 - We shouldn't commit vendor folder to repository, we have Godeps to handle it.
 - Add `1.10` Go version to travis.

#### How it works

```
go build -o namebeta
./namebeta
./namebeta alex.com
```

#### main.go

Move this code to utils.go:

```
if len(os.Args) == 1 {
	displayUsage()
	os.Exit(0)
}
```

Can be combined into:

```
func main() {
	query(parseArgs(os.Args))
}
```

Also these 3 vars can be wrapped into a struct, so later we can easily change struct and don't change code which is handling return values.

It's not clear to do Exit and return values in one func.

#### utils.go

Combine whoisQuery and domainQuery.

```
func getQueryResults(endpoint string, domain string, param map[string]string) []interface{} {
	var result []interface{}

	request := gorequest.New()
	_, body, _ := request.Post(endpoint).
		Type("form").
		Set("User-Agent", userAgent).
		Set("Refer", fmt.Sprintf(referURL, domain)).
		SendMap(param).End()

	if err := json.Unmarshal([]byte(body), &result); err != nil {
		color.Red(fmt.Sprintf("%s gailed to query %s endpoint. domain: %s \r\n", crossSymbol, endpoint, domain))
		os.Exit(1)
	}

	return result
}
```

Build Run!

