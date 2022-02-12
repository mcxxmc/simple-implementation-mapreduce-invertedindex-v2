package test

import (
	"simple-implementation-mapreduce-invertedindex-v2/indexing"
	"sort"
	"testing"
)

func TestPriorityQueue(t *testing.T) {
	data := [][]int{  // worker id: workload
		{1, 100}, {2, 50}, {3, 80}, {4, 200}, {5, 10}, {6, 150}, {7, 60}, {8, 500},
	}
	sortedData := [][]int{  // worker id: workload
		{1, 100}, {2, 50}, {3, 80}, {4, 200}, {5, 10}, {6, 150}, {7, 60}, {8, 500},
	}
	sort.Slice(sortedData, func(i, j int) bool {
		return sortedData[i][1] < sortedData[j][1]
	})

	pq := indexing.NewPriorityQueue(len(data))
	for _, rec := range data {
		pq.Add(rec[0], rec[1])
	}
	if pq.Num() != len(data) {
		t.Error("wrong number of elements: ", pq.Num(), ", expected ", len(data))
	}

	index := 0
	for pq.Num() > 0 {
		record := pq.ExtractMin()
		expected := sortedData[index]
		if record.WorkerId != expected[0] || record.Workload != expected[1] {
			t.Error("wrong worker id or workload; expected ", expected[0], " ", expected[1], ", but got: ",
				record.WorkerId, " ", record.Workload)
		}
		index ++
	}

	return
}
