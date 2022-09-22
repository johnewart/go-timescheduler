# go-timescheduler

## Introduction

This is a library for efficiently storing and retrieving time-based events. 
Events are scheduled by being placed into buckets based on the time they are due; 
this allows for efficient retrieval of events that are due at a given time by reducing 
the search space to only the currently-due bucket. 

Interacting with the scheduler will dynamically re-distribute events as they are due, handling 
overdue and due-beyond-scope events as well (e.g. events that are due in 10 years go into the last bucket 
and events that are overdue are placed into the first bucket).

## Example usage

```go 
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
```

Running this example will print out the following:

```bash

DUMPING!
TimespanBucket: 2022-09-22 11:34:11.51402 -0700 PDT m=+0.000112834 -> 2022-09-22 11:34:13.51402 -0700 PDT m=+2.000112834 (0)
TimespanBucket: 2022-09-22 11:34:13.514031 -0700 PDT m=+2.000123292 -> 2022-09-22 11:34:15.514031 -0700 PDT m=+4.000123292 (0)
TimespanBucket: 2022-09-22 11:34:15.514031 -0700 PDT m=+4.000123501 -> 2022-09-22 11:34:17.514031 -0700 PDT m=+6.000123501 (1)
 * birthday! @ %!s(func() time.Time=0x1026d1ea0)
TimespanBucket: 2022-09-22 11:34:17.514033 -0700 PDT m=+6.000126084 -> 2022-09-22 11:34:19.514033 -0700 PDT m=+8.000126084 (0)
TimespanBucket: 2022-09-22 11:34:19.514033 -0700 PDT m=+8.000126167 -> 2022-09-22 11:34:21.514033 -0700 PDT m=+10.000126167 (0)
TimespanBucket: 2022-09-22 11:34:21.514034 -0700 PDT m=+10.000126376 -> 2022-09-22 11:34:23.514034 -0700 PDT m=+12.000126376 (0)
TimespanBucket: 2022-09-22 11:34:23.514034 -0700 PDT m=+12.000126501 -> 2022-09-22 11:34:25.514034 -0700 PDT m=+14.000126501 (0)
TimespanBucket: 2022-09-22 11:34:25.514034 -0700 PDT m=+14.000126584 -> 2022-09-22 11:34:27.514034 -0700 PDT m=+16.000126584 (0)
TimespanBucket: 2022-09-22 11:34:27.514034 -0700 PDT m=+16.000126667 -> 2022-09-22 11:34:29.514034 -0700 PDT m=+18.000126667 (0)
TimespanBucket: 2022-09-22 11:34:29.514034 -0700 PDT m=+18.000127167 -> 2022-09-22 11:34:31.514034 -0700 PDT m=+20.000127167 (0)
TimespanBucket: 2022-09-22 11:34:31.514034 -0700 PDT m=+20.000127167 -> 2022-09-22 11:34:33.514034 -0700 PDT m=+22.000127167 (0)
TimespanBucket: 2022-09-22 11:34:33.514034 -0700 PDT m=+22.000127167 -> 2022-09-22 11:34:35.514034 -0700 PDT m=+24.000127167 (0)
```

And several seconds later it will print out the text below; notice that the buckets that are in the past have been removed and 
new buckets have been added to the end of the schedule. 

```bash 
DUMPING!
TimespanBucket: 2022-09-22 11:34:15.514031 -0700 PDT m=+4.000123501 -> 2022-09-22 11:34:17.514031 -0700 PDT m=+6.000123501 (1)
 * birthday! @ %!s(func() time.Time=0x1026d1ea0)
TimespanBucket: 2022-09-22 11:34:17.514033 -0700 PDT m=+6.000126084 -> 2022-09-22 11:34:19.514033 -0700 PDT m=+8.000126084 (0)
TimespanBucket: 2022-09-22 11:34:19.514033 -0700 PDT m=+8.000126167 -> 2022-09-22 11:34:21.514033 -0700 PDT m=+10.000126167 (0)
TimespanBucket: 2022-09-22 11:34:21.514034 -0700 PDT m=+10.000126376 -> 2022-09-22 11:34:23.514034 -0700 PDT m=+12.000126376 (0)
TimespanBucket: 2022-09-22 11:34:23.514034 -0700 PDT m=+12.000126501 -> 2022-09-22 11:34:25.514034 -0700 PDT m=+14.000126501 (0)
TimespanBucket: 2022-09-22 11:34:25.514034 -0700 PDT m=+14.000126584 -> 2022-09-22 11:34:27.514034 -0700 PDT m=+16.000126584 (0)
TimespanBucket: 2022-09-22 11:34:27.514034 -0700 PDT m=+16.000126667 -> 2022-09-22 11:34:29.514034 -0700 PDT m=+18.000126667 (0)
TimespanBucket: 2022-09-22 11:34:29.514034 -0700 PDT m=+18.000127167 -> 2022-09-22 11:34:31.514034 -0700 PDT m=+20.000127167 (0)
TimespanBucket: 2022-09-22 11:34:31.514034 -0700 PDT m=+20.000127167 -> 2022-09-22 11:34:33.514034 -0700 PDT m=+22.000127167 (0)
TimespanBucket: 2022-09-22 11:34:33.514034 -0700 PDT m=+22.000127167 -> 2022-09-22 11:34:35.514034 -0700 PDT m=+24.000127167 (0)
TimespanBucket: 2022-09-22 11:34:35.514034 -0700 PDT m=+24.000127167 -> 2022-09-22 11:34:37.514034 -0700 PDT m=+26.000127167 (0)
TimespanBucket: 2022-09-22 11:34:37.514034 -0700 PDT m=+26.000127167 -> 2022-09-22 11:34:39.514034 -0700 PDT m=+28.000127167 (0)
TimespanBucket: 2022-09-22 11:34:39.514034 -0700 PDT m=+28.000127167 -> 2022-09-22 11:34:41.514034 -0700 PDT m=+30.000127167 (0)

```

If you wait a little longer without removing the item you will see that when the first bucket has been rolled off, the scheduler will automatically
place any unprocessed events into what is now the current first bucket.

```bash 
TimespanBucket: 2022-09-22 11:35:57.907857 -0700 PDT m=+10.000125084 -> 2022-09-22 11:35:59.907857 -0700 PDT m=+12.000125084 (1)
 * birthday! @ %!s(func() time.Time=0x102555ea0)
```