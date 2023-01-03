package main

import (
	"fmt"
	"time"
)

func action() {
	fmt.Println("FUck time is over")
}

var myTimer *time.Timer

func timerFunc() {
	time_in_seconds := 30
	myTimer = NewTimer(time_in_seconds, func() {
		fmt.Printf("Congratulations! Your %d second timer finished.", time_in_seconds)
	})
	defer myTimer.Stop()
}

func NewTimer(seconds int, action func()) *time.Timer {
	myTimer = time.NewTimer(time.Duration(seconds) * time.Second)
	for {
		select {
		case <-myTimer.C:
			action()
		}
	}
	return myTimer
}
