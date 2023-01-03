package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var sigKillCounter = 0

func createError() {
	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel,
		syscall.SIGINT,
		syscall.SIGQUIT)
	exit_chan := make(chan int)
	go func() {
		for {
			s := <-signalChanel
			switch s {
			// kill -SIGINT XXXX или Ctrl+c  [XXXX - идентификатор процесса для программы]
			case syscall.SIGINT:
				sigKillCounter += 1
				fmt.Println("Signal interrupt triggered.")
				if sigKillCounter >= 3 {
					exit_chan <- 1
				}
			default:
				fmt.Println("Unknown signal.")
				exit_chan <- 1
			}
		}
	}()
	exitCode := <-exit_chan
	os.Exit(exitCode)
}
