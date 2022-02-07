package test

import (
	"simple-implementation-mapreduce-invertedindex-v2/indexing"
	"testing"
)

func TestQueue(t *testing.T) {
	queueSize := 10
	q := indexing.NewQueue(queueSize)
	if _, b := q.Pop(); b || q.Capacity() != queueSize || q.NumberOfItems() != 0 {
		t.Error("t1")
		return
	}

	in := 0  // the next input
	items := 0
	for in < 3 {
		q.Push(in)
		in ++
		items ++
	}
	out := 0  // the expected output
	if v, b := q.Peak(); !b || v != out || q.Capacity() != queueSize || q.NumberOfItems() != items {
		t.Error("t2")
		return
	}

	v, b := q.Pop()
	items --
	if v != out || !b || q.NumberOfItems() != items {
		t.Error("t3")
		return
	}
	out ++

	for items < queueSize {
		q.Push(in)
		in ++
		items ++
	}
	if q.NumberOfItems() != queueSize {
		t.Error("t4")
		return
	}

	b = q.Push(in)  // this should return false as the queue is filled full
	if b {
		t.Error("t5")
		return
	}

	for out < in - 2 {
		v, b = q.Pop()
		if v != out || !b {
			t.Error("t6")
			return
		}
		out ++
		items --
	}
	if q.NumberOfItems() != in - out {
		t.Error("t7")
		return
	}

	for items < q.Capacity() {
		q.Push(in)
		in ++
		items ++
	}
	for out < in {
		if v, b = q.Pop(); v != out || !b {
			t.Error("t8")
			return
		}
		out ++
	}

	if _, b = q.Pop(); b {
		t.Error("t9")
		return
	}
}
