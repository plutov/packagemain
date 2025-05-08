package main

import (
	json "encoding/json/v2"
	"fmt"
	"log"
	"time"
)

type Target struct {
	Msg  string    `json:"msg"`
	Time time.Time `json:"time,format:'2006-01-02'"`
}

var in = `{
	"msg": "hello1",
	"Msg": "hello3",
	"time": "2025-05-06"
}`

func main() {
	out := Target{}
	if err := json.Unmarshal([]byte(in), &out); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("out: %+v\n", out)
}
