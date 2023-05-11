package main

import (
	"time"

	douBot "github.com/brongh/go-telegram-financial-bot/douBot"
)

func main(){
	go douBot.StartBot()

	for {
		time.Sleep(60 * time.Second)
	}
}