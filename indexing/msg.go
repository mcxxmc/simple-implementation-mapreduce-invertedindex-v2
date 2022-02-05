package indexing

// Msg the message for communication between collector and manager.
type Msg struct {
	data interface{}    // the data; the real type depends on the typ
	typ int  		    // the type of the msg
	id int  			// the id of the sender
}

func newMsg(data interface{}, typ int, id int) *Msg {
	return &Msg{data: data, typ: typ, id: id}
}

// sends msg from collector to manager to show that it is idle.
func newMsgCollectorIdle(id int) *Msg {
	return &Msg{typ: msgCollectorIdle, id: id}
}

// sends msg from collector to manager to show that it is busy.
func newMsgCollectorBusy(id int) *Msg {
	return &Msg{typ: msgCollectorBusy, id: id}
}
