package main

import (
	"time"

	"github.com/vote_bot/modules"
	"github.com/vote_bot/util"
)

func main() {
	util.Init()
	go modules.Telegram()
	time.Sleep(2 * time.Second)
	for {
		modules.Proposals()
		time.Sleep(5 * time.Minute)
	}
}
