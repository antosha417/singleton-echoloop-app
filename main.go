package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	SecretKey   = "very secret"
	MaxTextSize = 20
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

type Data struct {
	Text   string `json:"text"`
	Secret string `json:"secret"`
}

func getNewText(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var data Data
	err := decoder.Decode(&data)
	panicOnError(err)
	log.Println(data)

	if data.Secret == SecretKey {
		if len(data.Text) > MaxTextSize {
			// TODO response with error
		}
		log.Printf("got new text:  <%s>\n", data.Text)
		// TODO add text to global array
	}
	_, err = fmt.Fprintln(w, "Some response text!")
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
		log.Print("Response: ", bodyString)
	}
}

func cleanArgs(text *string, port *int) (err error) {
	if len(*text) > MaxTextSize {
		err = fmt.Errorf("text is too long, max size is %d", MaxTextSize)
	}

	if *port < 0 || *port > 1<<16-1 {
		err = fmt.Errorf("invalid port %d", *port)
	}
	if err != nil {
		flag.Usage()
	}
	return
}

func main() {

	text := flag.String("text", "some string", "text to echo every second")
	port := flag.Int("port", 8080, "port to listen new text")
	url := flag.String("url", "http://localhost", "url to send text")
	flag.Parse()
	err := cleanArgs(text, port)
	panicOnError(err)

	http.HandleFunc("/", getNewText)
	log.Printf("Attempt to start server at :%d\n", *port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)

	if err != nil {
		log.Println(err.Error())
		log.Printf("can't listen on %d port, the first instance must be listening", *port)
		sendTextToFirstInstance(*text, fmt.Sprintf("%s:%d", *url, *port))
		log.Printf("echo loop for <%s> finished!", *text)
		return
	}
	// TODO echo loop every second
}
