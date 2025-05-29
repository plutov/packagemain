package survey

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

func SendSurvey(s *Survey, w io.Writer) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	return err
}

func example() {
	file, err := os.Open("survey.out")
	if err != nil {
		log.Fatal(err)
	}

	s := &Survey{
		Title: "My Survey",
	}
	SendSurvey(s, file)
}
