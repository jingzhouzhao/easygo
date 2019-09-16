package main

import (
	_ "easygo/conf"
	"easygo/server"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	server.Start()
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	for {
		s := <-exit
		fmt.Println("exit:", s)
		if err := server.Close(); err != nil {
			panic(err)
		}
		fmt.Println("exited!")
		return
	}
}
