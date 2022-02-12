package main

import (
	"simple-implementation-mapreduce-invertedindex-v2/common"
	"simple-implementation-mapreduce-invertedindex-v2/indexing"
	"strconv"
)

const numOfWorkers = 8

// modify this to change the initial patches of jobs
func createJobs() []string {
	r := make([]string, 61)
	for i := 0; i < len(r); i ++ {
		r[i] = common.PapDividedPathPrefixMain + strconv.Itoa(i + 1) + common.TxtAppendix
	}
	return r
}

func main() {
	//preprocess.SplitOriginalTxt()
	jobs := createJobs()
	manager := indexing.NewManager(numOfWorkers)
	manager.Run(jobs)
}
