package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	SecretKey   = "very secret"
	MaxTextSize = 20
	MaxTexts    = 10
)

type Texts struct {
	sync.Mutex
	Texts []string
}

func NewTexts() Texts {
	return Texts{Texts: make([]string, 0, MaxTexts)}
}

func (t *Texts) Append(text string) {
	t.Lock()
	t.Texts = append(t.Texts, text)
	t.Unlock()
}

func (t *Texts) count() int {
	t.Lock()
	defer t.Unlock()
	return len(t.Texts)
}

var AllTexts Texts

func echoLoop() {

	ticker := time.Tick(5 * time.Second)
	i := 0
	for range ticker {
		AllTexts.Lock()
		i = (i + 1) % len(AllTexts.Texts)
		val := AllTexts.Texts[i]
		AllTexts.Unlock()

		fmt.Println(val)
	}
}
