package indexing

import "simple-implementation-mapreduce-invertedindex-v2/common"

// Msg the message for communication between collector and manager.
type Msg struct {
	Data interface{} // the data; the real type depends on the Typ
	Typ  int         // the type of the msg
	Id   int         // the id of the sender
}

//func NewMsg(data interface{}, Typ int, id int) *Msg {
//	return &Msg{data: data, Typ: Typ, Id: id}
//}

// NewMsgCountFreq order from manager to collector to count the word frequency in a certain file.
func NewMsgCountFreq(filename string, id int) *Msg {
	return &Msg{Data: filename, Typ: MsgCountFreq, Id: id}
}

// NewMsgCombineFreq order from manager to collector to combine records.
func NewMsgCombineFreq(records common.NativeRecords, id int) *Msg {
	return &Msg{Data: records, Typ: MsgCombineFreq, Id: id}
}

// NewMsgDismissWorker order from manager to collector to dismiss it.
func NewMsgDismissWorker(id int) *Msg {
	return &Msg{Typ: MsgDismissWorker, Id: id}
}

// NewMsgDeliverData order from manager to collector to deliver its data (records).
func NewMsgDeliverData(id int) *Msg {
	return &Msg{Typ: MsgDeliverData, Id: id}
}

// NewMsgClearData order from manager to collector to clear its data (records).
func NewMsgClearData(id int) *Msg {
	return &Msg{Typ: MsgClearData, Id: id}
}

// NewMsgSortSave2Disk order from manager to collector to sort and save its records on the disk.
func NewMsgSortSave2Disk(savePath string, id int) *Msg {
	return &Msg{Data: savePath, Typ: MsgSortAndSave2Disk, Id: id}
}

// NewMsgCollectorIdle sends msg from collector to manager to show that it is idle.
func NewMsgCollectorIdle(id int) *Msg {
	return &Msg{Typ: MsgCollectorIdle, Id: id}
}

// NewMsgCollectorDelivery sends msg from collector to manager to deliver data (records).
func NewMsgCollectorDelivery(records common.NativeRecords, id int) *Msg {
	return &Msg{Data: records, Typ: MsgCollectorDelivery, Id: id}
}
