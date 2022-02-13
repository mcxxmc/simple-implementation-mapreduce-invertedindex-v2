package common

import "strconv"

// Record the record to be stored on the disk in a JSON file.
type Record struct {
	Source string `json:"source"`
	Freq int `json:"freq"`
}

func (rec *Record) String() string {
	return "[ source: " + rec.Source + ", freq: " + strconv.Itoa(rec.Freq) + " ]\n"
}

func NewRecord(source string, freq int) *Record {
	return &Record{Source: source, Freq: freq}
}

// Records the records to be stored on the disk as a JSON file.
//
// in the form word: []*Record
type Records map[string][]*Record

func NewRecords() Records {
	return Records{}
}

// NativeRecord the record to be stored by the living collectors.
//
// in the form filename: frequency
type NativeRecord map[string]int

func NewNativeRecord() NativeRecord {
	return NativeRecord{}
}

// NativeRecords the records to be stored by the living collectors.
//
// in the form word: NativeRecord
type NativeRecords map[string]NativeRecord

func NewNativeRecords() NativeRecords {
	return NativeRecords{}
}
