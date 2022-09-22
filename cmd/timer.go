package main

import (
	"context"
	"fmt"
	"github.com/johnewart/go-timescheduler/schedule"
	"time"
)

type Birthday struct {
	schedule.Schedulable
}

func (b Birthday) DueTime() time.Time {
	return time.Now().Add(time.Second * 5)
}

func (b Birthday) Id() string {
	return "birthday!"
}

func main() {
	ctx := context.Background()
	scheduler := schedule.NewScheduler[Birthday](ctx, time.Second*2, 10)
	scheduler.AddReminder(Birthday{})
	for {
		fmt.Println("DUMPING!")
		scheduler.Dump()
		fmt.Println("----")
		time.Sleep(5 * time.Second)
	}
}
