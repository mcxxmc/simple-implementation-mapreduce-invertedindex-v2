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
	id       int
	idle 	 bool
	alive 	 bool
}

// NewCollector returns the pointer to a new Collector.
func NewCollector(receiver chan *Msg, sender chan *Msg, id int) *Collector {
	return &Collector{
		records: common.NewNativeRecords(),
		receiver: receiver, sender: sender, id: id,
		idle: true, alive: false,
	}
}

// Run initializes and runs the collector.
func (c *Collector) Run() {
	c.alive = true
	for c.alive {
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
				c.setBusy()
				c.countFreq(job)
				c.setIdle()

			case MsgCombineFreq:
				records, err := msg.Data.(common.NativeRecords)
				if !err {
					fmt.Println("wrong data type")
					break
				}
				fmt.Println("collector id: ", c.id, " starts combining frequency")
				c.setBusy()
				c.combineFreq(records)
				c.setIdle()

			case MsgCurrentStatus:
				// only make responses if the collector is idle
				if c.idle {
					c.sender <- NewMsgCollectorIdle(c.id)
				}

			case MsgDismissWorker:
				fmt.Println("collector dismissed, id: ", c.id)
				c.idle = true
				c.alive = false

			case MsgDeliverData:
				fmt.Println("collector id: ", c.id, " starts delivering data")
				c.setBusy()
				newMsg := NewMsgCollectorDelivery(c.records, c.id)
				c.sender <- newMsg
				c.setIdle()

			case MsgClearData:
				fmt.Println("collector id: ", c.id, " starts cleaning data")
				c.setBusy()
				c.records = common.NewNativeRecords()
				c.setIdle()

			case MsgSortAndSave2Disk:
				savePath, err := msg.Data.(string)
				if !err {
					fmt.Println("wrong data type")
					break
				}
				fmt.Println("collector id: ", c.id, " starts sorting and saving data")
				c.setBusy()
				c.sortSave(savePath)
				c.setIdle()

			default:
				fmt.Println("unknown message type: ", msg.Typ)
			}

		case <- time.After(2 * time.Second):
			continue
		}
	}
}

func (c *Collector) setBusy() {
	c.idle = false
	//c.sender <- newMsgCollectorBusy(c.id)
}

func (c *Collector) setIdle() {
	c.sender <- NewMsgCollectorIdle(c.id)
	c.idle = true
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
		case b >= '0' && b <= '9' || b == '-' || b >= 'a' && b <= 'z':
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
