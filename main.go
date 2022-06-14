package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"szpt/szpt"
)

func main() {
	fmt.Println("service start")
	c := cron.New()
	c.AddFunc("1 0 * * *", func() {
		szpt.MultiReport()
		fmt.Printf("time now")
	})
	c.Start()
	select {}
}
