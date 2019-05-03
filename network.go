package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Data struct {
	Text   string `json:"text"`
	Secret string `json:"secret"`
}

func cleanNewText(data *Data) error {
	msg := ""
	if data.Secret != SecretKey {
		msg = "secrets don't match"
	}
	if len(data.Text) > MaxTextSize {
		msg = "text is too long"
	}
	if AllTexts.count() >= MaxTexts {
		msg = "there are already too many texts to echo"
	}
	if msg != "" {
		return fmt.Errorf("error: %s", msg)
	}
	return nil
}

func waitForHttpServerToStart(url string) {
	log.Println("waiting for http server to start")
	for {
		time.Sleep(time.Second)

		log.Println("checking if started...")
		resp, err := http.Get(url + "/ping")
		if err != nil {
			log.Println("failed:", err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			log.Println("not ok:", resp.StatusCode)
			panicOnError(resp.Body.Close())
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		if stringBody := string(body); stringBody != "pong" {
			log.Printf("wrong ans:<%s>", stringBody)
			panicOnError(resp.Body.Close())
			continue
		}
		break
	}
	log.Println("server up and running!")
}

func pong(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "pong")
	panicOnError(err)
}

func getNewText(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var data Data
	err := decoder.Decode(&data)
	panicOnError(err)

	respMessage := "ok!"
	if err = cleanNewText(&data); err == nil {
		AllTexts.Append(data.Text)
	} else {
		respMessage = err.Error()
	}
	_, err = fmt.Fprint(w, respMessage)
	panicOnError(err)
}

func sendTextToFirstInstance(text string, url string) {
	log.Printf("sending <%s> to the first instance (%s)", text, url)
	data := Data{Text: text, Secret: SecretKey}
	marshaledData, err := json.Marshal(data)
	panicOnError(err)
	resp, err := http.Post(url, "application/json", bytes.NewReader(marshaledData))
	panicOnError(err)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		panicOnError(err)
		bodyString := string(bodyBytes)
		log.Print("response: ", bodyString)
	}
}
