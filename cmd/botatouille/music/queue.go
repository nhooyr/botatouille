package music

import (
	"sync"
)

// TODO turn it into its own process maybe?
type queue struct {
	playing   *video
	head      *node
	startChan chan<- *video
	sync.Mutex
}

func newQueue(startChan chan<- *video) *queue {
	return &queue{startChan: startChan}
}

// TODO maybe use value only?
type video struct {
	id    string
	title string
}

func (q *queue) getPlaying() *video {
	q.Lock()
	defer q.Unlock()
	return q.playing
}

func (q *queue) append(v *video) {
	q.Lock()
	if q.playing == nil {
		q.playing = v
		q.startChan <- v
		// TODO start playing somehow?
	} else if q.head == nil {
		q.head = &node{video: v}
	} else {
		q.head.append(v)
	}
	q.Unlock()
}

func (q *queue) next() *video {
	q.Lock()
	defer q.Unlock()
	if q.head == nil {
		q.playing = nil
		return nil
	}
	q.playing = q.head.video
	v := q.playing
	q.head = q.head.next
	return v
}

func (q *queue) playlist() []*video {
	q.Lock()
	defer q.Unlock()
	if q.head == nil {
		return []*video{}
	}
	return q.head.playlist(nil)
}

type node struct {
	next  *node
	video *video
}

func (n *node) append(v *video) {
	if n.next == nil {
		n.next = &node{video: v}
		return
	}
	n.next.append(v)
}

func (n *node) playlist(p []*video) []*video {
	p = append(p, n.video)
	if n.next == nil {
		return p
	}
	return n.next.playlist(p)
}
