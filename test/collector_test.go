package test

import (
	"simple-implementation-mapreduce-invertedindex-v2/common"
	"simple-implementation-mapreduce-invertedindex-v2/indexing"
	"testing"
)

func TestSingleCollector(t *testing.T) {
	testFilename := common.PapDividedPathPrefix + "1" + common.TxtAppendix
	savePath := "testSingleCollector.json"
	receiver := make(chan *indexing.Msg, indexing.CollectorChanCapacity)
	sender := make(chan *indexing.Msg, indexing.ManagerChanCapacity)
	collector := indexing.NewCollector(receiver, sender, 1)
	go collector.Run()
	receiver <- indexing.NewMsgCountFreq(testFilename, 0)  // 1
	receiver <- indexing.NewMsgSortSave2Disk(savePath, 0)  // 1
	<- sender
}

func Test2Collectors(t *testing.T) {
	testFilename1 := common.PapDividedPathPrefix + "1" + common.TxtAppendix
	testFilename2 := common.PapDividedPathPrefix + "2" + common.TxtAppendix
	savePath := "test2Collectors.json"
	receiver1 := make(chan *indexing.Msg, indexing.CollectorChanCapacity)
	receiver2 := make(chan *indexing.Msg, indexing.CollectorChanCapacity)
	sender := make(chan *indexing.Msg, indexing.ManagerChanCapacity)
	collector1 := indexing.NewCollector(receiver1, sender, 1)
	collector2 := indexing.NewCollector(receiver2, sender, 2)

	go collector1.Run()
	go collector2.Run()

	receiver1 <- indexing.NewMsgCountFreq(testFilename1, 0)  // 1
	receiver2 <- indexing.NewMsgCountFreq(testFilename2, 0)  // 1
	receiver2 <- indexing.NewMsgDeliverData(0)  // 2

	msg := <- sender
	for msg.Typ != indexing.MsgCollectorDelivery {
		t.Error("wrong msg type")
	}

	receiver1 <- indexing.NewMsgCombineFreq(msg.Data.(common.NativeRecords), 0)  // 1
	receiver2 <- indexing.NewMsgClearData(0)  // 1
	receiver1 <- indexing.NewMsgSortSave2Disk(savePath, 0)  // 1
	<- sender
}
