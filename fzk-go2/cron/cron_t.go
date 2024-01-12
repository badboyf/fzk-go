package cron

import (
	"flag"
	"fmt"
	"github.com/robfig/cron"
	"log"
	"time"
)

func CronT1() {
	spec := flag.String("s", "0 * * * * ?", "cron spec")
	next := flag.Uint("i", 5, "cron next times")
	flag.Parse()

	p, err := cron.Parse(*spec)
	if err != nil {
		log.Fatal(err)
	}
	t := time.Now()
	for i := int(*next); i > 0; i-- {
		t = p.Next(t)
		fmt.Println(t.String())
	}
}
