package indexing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"simple-implementation-mapreduce-invertedindex-v2/common"
	"sort"
	"time"
)

// Collector the worker for creating inverted index in a concurrent manner.
//
// Please use NewCollector() as the constructor.
type Collector struct {
	records  common.NativeRecords
	receiver chan *Msg
	sender 	 chan *Msg
	id     int
	online bool
}

// NewCollector returns the pointer to a new Collector.
func NewCollector(receiver chan *Msg, sender chan *Msg, id int) *Collector {
	return &Collector{
		records: common.NewNativeRecords(),
		receiver: receiver, sender: sender, id: id,
		online: false,
	}
}

// Run initializes and runs the collector.
func (c *Collector) Run() {
	c.online = true
	for c.online {
		select {
		case msg := <- c.receiver:
			switch msg.Typ {
			case MsgCountFreq:
				job, err := msg.Data.(string)
				if !err {
					fmt.Println("wrong data type")
					break
				}
				fmt.Println("collector id: ", c.id, " starts counting frequency for: ", job)
				c.countFreq(job)
				fmt.Println("collector id: ", c.id, " has finished counting frequency for: ", job)

			case MsgCombineFreq:
				records, err := msg.Data.(common.NativeRecords)
				if !err {
					fmt.Println("wrong data type")
					break
				}
				fmt.Println("collector id: ", c.id, " starts combining frequency")
				c.combineFreq(records)
				fmt.Println("collector id: ", c.id, " has finished combining frequency")

			case MsgDismissWorker:
				fmt.Println("collector dismissed, id: ", c.id)
				c.online = false

			case MsgDeliverData:
				fmt.Println("collector id: ", c.id, " starts delivering data")
				c.sender <- NewMsgCollectorDelivery(c.records, c.id)
				fmt.Println("collector id: ", c.id, " has finished delivering data")

			case MsgClearData:
				fmt.Println("collector id: ", c.id, " starts cleaning data")
				c.records = common.NewNativeRecords()
				fmt.Println("collector id: ", c.id, " has finished cleaning data")

			case MsgSortAndSave2Disk:
				savePath, err := msg.Data.(string)
				if !err {
					fmt.Println("wrong data type")
					break
				}
				fmt.Println("collector id: ", c.id, " starts sorting and saving data")
				c.sortSave(savePath)
				fmt.Println("collector id: ", c.id, " has finished sorting and saving data")
				c.sender <- NewMsgCollectorIOCompleted(c.id)

			default:
				fmt.Println("unknown message type: ", msg.Typ)
			}

		case <- time.After(2 * time.Second):
			continue
		}
	}
}

// counts the word frequency. Currently, it assumes there is no hash collision.
func (c *Collector) countFreq(path string) {
	byt, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	word := ""

	for _, b := range byt {
		switch {
		case b >= '0' && b <= '9' || b >= 'a' && b <= 'z':
			word += string(b)
		case b >= 'A' && b <= 'Z':
			word += string(b + 'a' - 'A')
		default:
			if word != "" {
				if wordRecord, exist := c.records[word]; exist {
					wordRecord[path] ++
				} else {
					c.records[word] = common.NewNativeRecord()
					c.records[word][path] = 1
				}
				word = ""
			}
		}
	}
	return
}

// combines the word frequency from another collector.
func (c *Collector) combineFreq(records common.NativeRecords) {
	for word, wordRecord := range records {
		if ownWordRecord, exist := c.records[word]; exist {
			for filename, freq := range wordRecord {
				ownWordRecord[filename] += freq
			}
		} else {
			c.records[word] = wordRecord
		}
	}
}

// converts, sorts and saves the records to the disk.
func (c *Collector) sortSave(savePath string) {
	records := common.NewRecords()
	for word, wordRecord := range c.records {
		tmp := make([]*common.Record, len(wordRecord))
		index := 0
		for filename, freq := range wordRecord {
			tmp[index] = common.NewRecord(filename, freq)
			index ++
		}
		sort.Slice(tmp, func(i, j int) bool {return tmp[i].Freq > tmp[j].Freq})
		records[word] = tmp
	}

	b, err := json.Marshal(records)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ioutil.WriteFile(savePath, b, 0644)
	if err != nil {
		fmt.Println(err)
	}
}
