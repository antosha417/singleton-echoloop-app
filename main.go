package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
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
	fullUrl := fmt.Sprintf("%s:%d", *url, *port)

	AllTexts = NewTexts()
	AllTexts.Texts = append(AllTexts.Texts, *text)
	go func() { waitForHttpServerToStart(fullUrl); echoLoop() }()

	http.HandleFunc("/ping", pong)
	http.HandleFunc("/", getNewText)
	log.Printf("attempt to start server at :%d\n", *port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)

	if err != nil {
		log.Println(err.Error())
		log.Printf("can't listen on %d port, the first instance must be listening", *port)
		sendTextToFirstInstance(*text, fullUrl)
		log.Printf("echo loop for <%s> finished!", *text)
		return
	}
}
