package main

import (
	"bytes"

	"github.com/utahta/go-linenotify"
)

func main() {
	token := "" // EDIT THIS

	c := linenotify.New(linenotify.WithToken(token))
	c.Notify("hello world", "", "", nil)

	c = linenotify.New()
	c.SetToken(token)
	c.Notify("hello world", "http://localhost/thumb.jpg", "http://localhost/full.jpg", nil)
	c.Notify("hello world", "", "", bytes.NewReader([]byte("image bytes")))
}
