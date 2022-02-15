package search

import (
	"simple-implementation-mapreduce-invertedindex-v2/common"
	"sort"
	"strconv"
	"strings"
)

// the summary for a single file.
type singleFileSummary struct {
	filename   string   // the filename
	relevantTo []string // all the keywords relevant to the file
	totalFreq  int      // the higher the score, the more relevant the file is
}

// updates a record. Note that it does not check if the filename is the same.
func (sfs *singleFileSummary) update(word string, rec *common.Record) {
	sfs.totalFreq += rec.Freq
	sfs.relevantTo = append(sfs.relevantTo, word)
}

func (sfs *singleFileSummary) String() string {
	return sfs.filename + "\n\t score: " + strconv.Itoa(sfs.totalFreq) +
		", with words: " + strings.Join(sfs.relevantTo, ", ")
}

// Summary the Summary object is used for the searching process when searching for multiple words;
//
// it contains the summary for all relevant files, in the form: filename: singleFileSummary.
//
// a file will only be included a summary when it contains the wanted words.
type Summary struct {
	S 		map[string]*singleFileSummary
}

// NewSummary returns a pointer to a new Summary object.
func NewSummary() *Summary {
	return &Summary{S: make(map[string]*singleFileSummary)}
}

// Add takes in a record.
func (s *Summary) Add(word string, rec *common.Record) {
	if sfs, exist := s.S[rec.Source]; exist {
		sfs.update(word, rec)
	} else {
		s.S[rec.Source] = &singleFileSummary{filename: rec.Source, relevantTo: []string{word}, totalFreq: rec.Freq}
	}
}

func (s *Summary) Sort() []*singleFileSummary {
	r := make([]*singleFileSummary, len(s.S))
	index := 0
	for _, sfs := range s.S {
		r[index] = sfs
		index ++
	}
	sort.Slice(r, func(i, j int) bool {
		// currently, the ranking solely depends on the total frequency
		return r[i].totalFreq > r[j].totalFreq
		/*
		// for a better comparison, use this instead; you are also welcome to implement your own ranking algorithm
		return r[i].totalFreq * len(r[i].relevantTo) > r[j].totalFreq * len(r[j].relevantTo)
		 */
	})
	return r
}
