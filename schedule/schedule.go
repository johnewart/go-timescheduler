package schedule

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Schedulable interface {
	DueTime() time.Time
	Id() string
}

type TimespanBucket[T Schedulable] struct {
	startTime time.Time
	endTime   time.Time
	elements  []T
	lock      *sync.Mutex
}

func NewTimespanBucket[T Schedulable](startTime time.Time, endTime time.Time) *TimespanBucket[T] {
	return &TimespanBucket[T]{
		startTime: startTime,
		endTime:   endTime,
		elements:  make([]T, 0),
		lock:      &sync.Mutex{},
	}
}

func (t *TimespanBucket[Schedulable]) Contains(in time.Time) bool {
	return t.startTime.Before(in) && t.endTime.After(in)
}

func (t *TimespanBucket[T]) AddEntity(entity T) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.elements = append(t.elements, entity)
}

func (t *TimespanBucket[T]) Past() bool {
	return time.Now().After(t.endTime)
}

func (t *TimespanBucket[T]) String() string {
	return fmt.Sprintf("TimespanBucket: %s -> %s", t.startTime, t.endTime)
}

func (t *TimespanBucket[T]) Size() int {
	t.lock.Lock()
	defer t.lock.Unlock()
	return len(t.elements)
}

func (t *TimespanBucket[T]) IsAfter(dueTime time.Time) bool {
	return t.startTime.After(dueTime)
}

func (t *TimespanBucket[T]) IsBefore(dueTime time.Time) bool {
	return t.endTime.Before(dueTime)
}

type Scheduler[T Schedulable] struct {
	buckets   []*TimespanBucket[T]
	blockSize time.Duration
	numBlocks int
	ctx       context.Context
	mutex     *sync.Mutex
}

func NewScheduler[T Schedulable](ctx context.Context, blockSize time.Duration, numBlocks int) *Scheduler[T] {
	buckets := make([]*TimespanBucket[T], 0)

	for i := 0; i < numBlocks; i++ {
		startTime := time.Now().Add(time.Duration(i) * blockSize)
		endTime := startTime.Add(blockSize)
		buckets = append(buckets, NewTimespanBucket[T](startTime, endTime))
	}

	return &Scheduler[T]{
		ctx:       ctx,
		buckets:   buckets,
		blockSize: blockSize,
		numBlocks: numBlocks,
		mutex:     &sync.Mutex{},
	}
}

func (s *Scheduler[T]) update() {
	overdueItems := make([]T, 0)
	startIdx := 0

	for idx, bucket := range s.buckets {
		if !bucket.Past() {
			startIdx = idx
			break
		}
	}

	for i := 0; i < startIdx; i++ {
		if s.buckets[i].Size() > 0 {
			// pass
			overdueItems = append(overdueItems, s.buckets[0].elements...)
		}
	}

	s.buckets = s.buckets[startIdx:]

	currentEndTime := s.buckets[len(s.buckets)-1].endTime
	newBuckets := make([]*TimespanBucket[T], 0)
	for j := 0; j <= startIdx; j++ {
		newBuckets = append(newBuckets, NewTimespanBucket[T](currentEndTime, currentEndTime.Add(s.blockSize)))
		currentEndTime = currentEndTime.Add(s.blockSize)
	}

	s.buckets = append(s.buckets, newBuckets...)

	s.buckets[0].elements = append(s.buckets[0].elements, overdueItems...)

}

func (s *Scheduler[T]) AddReminder(entity T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.update()

	if s.buckets[0].IsAfter(entity.DueTime()) {
		// Overdue? Put it at the head of the queue
		s.buckets[0].AddEntity(entity)
		return
	}

	if s.buckets[len(s.buckets)-1].IsBefore(entity.DueTime()) {
		// Too far out? Shove it into the last bucket
		s.buckets[len(s.buckets)-1].AddEntity(entity)
		return
	}

	for _, bucket := range s.buckets {
		if bucket.Contains(entity.DueTime()) {
			bucket.AddEntity(entity)
			return
		}
	}

}

func (s *Scheduler[T]) Due() []T {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.update()

	dueItems := make([]T, 0)

	bucket := s.buckets[0]

	removeIdxs := make([]int, 0)
	for i, entity := range bucket.elements {
		if entity.DueTime().Before(time.Now()) {
			dueItems = append(dueItems, entity)
			removeIdxs = append(removeIdxs, i)
		}
	}

	for i, idx := range removeIdxs {
		// Do the truffle shuffle!
		realIdx := idx - i // we have removed i elements so we need to subtract that from the index
		bucket.elements = append(bucket.elements[:realIdx], bucket.elements[realIdx+1:]...)
	}

	return dueItems
}

func (s *Scheduler[T]) Dump() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.update()

	for _, bucket := range s.buckets {
		fmt.Printf("%s (%d)\n", bucket.String(), bucket.Size())
		for _, entity := range bucket.elements {
			fmt.Printf(" * %s @ %s\n", entity.Id(), entity.DueTime)
		}
	}
}
