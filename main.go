package main

import (
	"time"
)

func main() {
	go StartServer()
	time.Sleep(100 * time.Millisecond)
	Test()
}
