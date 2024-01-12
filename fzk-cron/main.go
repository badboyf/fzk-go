package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/robfig/cron"
	"time"
)

func main() {
	cronNext()
}

func cronNext() {
	layout := "2006-01-02 15:04:05"

	spec := flag.String("s", "0 * * * * ?", "cron spec")
	next := flag.Uint("i", 5, "cron next times")
	tm_addr := flag.String("t", "2023-01-01 00:00:00", "base time")
	flag.Parse()
	//fmt.Printf("spec %s\n", *spec)
	//fmt.Printf("next %s\n", *next)
	//fmt.Printf("tm %s\n", *tm_addr)

	tm := *tm_addr
	t := time.Now()
	var err error
	if tm != "" {
		t, err = time.Parse(layout, tm)
		if err != nil {
			d, _ := json.Marshal(map[string]interface{}{"code": 1, "msg": err.Error()})
			fmt.Printf("%s", string(d))
			return
		}
	}
	flag.Parse()

	p, err := cron.Parse(*spec)
	if err != nil {
		d, _ := json.Marshal(map[string]interface{}{"code": 1, "msg": err.Error()})
		fmt.Printf("%s", string(d))
		return
	}
	data := make([]string, 0)
	for i := int(*next); i > 0; i-- {
		t = p.Next(t)
		data = append(data, t.Format(layout))
	}
	d, err := json.Marshal(map[string]interface{}{
		"code": 0,
		"data": data,
	})
	fmt.Println(string(d))
}
