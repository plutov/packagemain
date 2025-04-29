package main

import (
	"encoding/json"
	jsonv2 "encoding/json/v2"
	"fmt"
	"time"
)

type Target struct {
	Msg  string    `json:"msg"`
	Time time.Time `json:"time,format:'2006-01-02'"`
}

var in = `{
	"msg": "hello1",
	"time": "2025-03-10"
}
`

func main() {
	out := Target{}
	if err := json.Unmarshal([]byte(in), &out); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("v1 res: %+v\n", out)

	out = Target{}
	if err := jsonv2.Unmarshal([]byte(in), &out); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("v2 res: %+v\n", out)
}
