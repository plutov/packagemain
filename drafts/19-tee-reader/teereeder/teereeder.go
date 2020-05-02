package teereeder

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"os"
)

func ValidateAndSaveCSVFile(csvFile io.Reader) error {
	var buf bytes.Buffer
	reader := bufio.NewReader(csvFile)
	reader.WriteTo(&buf)

	csvReader := csv.NewReader(csvFile)

	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if len(line[0]) == 0 {
			return errors.New("rows cannot be empty")
		}
	}

	w, err := os.OpenFile("test.csv", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, &buf)
	return err
}

func ValidateAndSaveCSVFile2(csvFile io.Reader) error {
	var buf bytes.Buffer
	tee := io.TeeReader(csvFile, &buf)

	csvReader := csv.NewReader(tee)

	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if len(line[0]) == 0 {
			return errors.New("rows cannot be empty")
		}
	}

	w, err := os.OpenFile("test.csv", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, &buf)
	return err
}
