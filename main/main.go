package main

import (
	"fmt"
	"simple-implementation-mapreduce-invertedindex-v2/common"
	"simple-implementation-mapreduce-invertedindex-v2/refine"
	"simple-implementation-mapreduce-invertedindex-v2/search"
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
	/*
	if !preprocess.SplitOriginalTxt() {
		return
	}
	 */

	/*
	jobs := createJobs()
	manager := indexing.NewManager(numOfWorkers)
	manager.Run(jobs)
	 */

	/* */
	jobs := createJobs()
	err := refine.Refine(jobs)
	if err != nil {
		fmt.Println(err)
	}


	engine, err := search.NewSearchEngine()
	if err != nil {
		fmt.Println(err)
		return
	}
	engine.Run()
}
