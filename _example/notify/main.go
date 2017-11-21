package main

import (
	"bytes"

	"github.com/utahta/go-linenotify"
)

func main() {
	token := "" // EDIT THIS

	c := linenotify.New()
	c.Notify(token, "hello world", "", "", nil)
	c.Notify(token, "hello world", "http://localhost/thumb.jpg", "http://localhost/full.jpg", nil)
	c.Notify(token, "hello world", "", "", bytes.NewReader([]byte("image bytes")))
}
