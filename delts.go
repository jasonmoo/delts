package delts

import (
	"container/list"
	"sync"
)

type (
	SortedDeltaStream struct {
		me sync.Mutex

		start, end        int64
		capturing, closed bool

		list     *list.List
		max_size int

		Output chan []int64
	}
)

func NewSortedDeltaStream(size int) *SortedDeltaStream {
	return &SortedDeltaStream{
		me:       sync.Mutex{},
		max_size: size,
		list:     list.New(),
		Output:   make(chan []int64, size),
	}
}

func (d *SortedDeltaStream) Add(n int64) {

	d.me.Lock()

	if d.closed {
		return
	}

	if !d.capturing {
		d.start = n
		d.end = n
		d.capturing = true
	} else if n == d.end+1 {
		d.end = n
		d.consume()
	} else {

		if d.list.Len() == 0 || d.list.Front().Value.(int64) > n {
			// first element or less than first element, push on front
			d.list.PushFront(n)
		} else if d.list.Back().Value.(int64) < n {
			// greater than last element, append to back
			d.list.PushBack(n)
		} else {
			// search for place in list backwards since probably
			cur := d.list.Back()
			for cur.Prev() != nil && cur.Prev().Value.(int64) > n {
				cur = cur.Prev()
			}
			d.list.InsertBefore(n, cur)

			// if buffer reaches capacity,
			// emit and consume next range
			if d.list.Len() > d.max_size {
				d.emit()
				d.capturing = false
				d.consume()
			}
		}

	}

	d.me.Unlock()

}

func (d *SortedDeltaStream) emit() {

	if d.start == d.end {
		d.Output <- []int64{d.start}
	} else {
		d.Output <- []int64{d.start, d.end}
	}

}

func (d *SortedDeltaStream) consume() {

	if d.list.Len() == 0 {
		return
	}

	if !d.capturing {
		front := d.list.Front()
		d.start = front.Value.(int64)
		d.end = d.start
		d.capturing = true
		d.list.Remove(front)
	}

	for d.list.Len() > 0 && d.list.Front().Value.(int64) == d.end+1 {
		d.end++
		d.list.Remove(d.list.Front())
	}

}

func (d *SortedDeltaStream) Emit() {

	d.me.Lock()
	d.emit()
	d.me.Unlock()

}

func (d *SortedDeltaStream) Close() {

	d.me.Lock()
	d.emit()
	d.capturing = false
	for d.list.Len() > 0 {
		d.consume()
		d.emit()
		d.capturing = false
	}
	close(d.Output)
	d.closed = true
	d.me.Unlock()

}
