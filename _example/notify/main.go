package main

import (
	"bytes"
	"context"

	"github.com/utahta/go-linenotify"
)

func main() {
	token := "" // EDIT THIS

	c := linenotify.NewClient()
	c.Notify(context.Background(), token, "hello world", "", "", nil)
	c.Notify(context.Background(), token, "hello world", "http://localhost/thumb.jpg", "http://localhost/full.jpg", nil)
	c.Notify(context.Background(), token, "hello world", "", "", bytes.NewReader([]byte("image bytes")))
}
