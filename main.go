package main

import (
	"net/http"
	"os"
	"time"

	douBot "github.com/brongh/go-telegram-financial-bot/douBot"
)

func main(){
	go douBot.StartBot()
	//
	// for {
	// 	time.Sleep(60 * time.Second)
	// }
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr: ":" + port,
		Handler: nil,
		ReadTimeout: 15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}