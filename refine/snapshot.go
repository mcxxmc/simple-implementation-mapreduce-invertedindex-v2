package refine

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"simple-implementation-mapreduce-invertedindex-v2/common"
)

const snapshotLength = 200	// the approximate length of a snapshot string
const skipLength = 6		// to skip some useless info in the head; adjust this number as needed

// Snapshots a collection of snapshots of the files.
type Snapshots map[string]string

func CreateSnapshots(filePaths []string) error {
	snapshots := make(map[string]string)
	for _, filePath := range filePaths {
		bytes, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}
		if len(bytes) > snapshotLength {
			bytes = bytes[skipLength: snapshotLength]
			toCut := 0
			for toCut < snapshotLength - skipLength{
				if bytes[toCut] != '\n' && bytes[toCut] != '\r' {
					break
				}
				toCut ++
			}
			bytes = bytes[toCut:]
		}
		snapshots[filePath] = string(bytes)
	}
	bytes, err := json.Marshal(snapshots)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(common.SnapshotsSavePath, bytes, 0644)
	return err
}
