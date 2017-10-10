package main

import "github.com/dandeliondeathray/nona/slackmessaging"

func main() {
	go slackmessaging.StartProbes(24689)

	chBlock := make(chan bool)
	<-chBlock
}
