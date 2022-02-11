package indexing

const CollectorChanCapacity = 64
const ManagerChanCapacity = 512

const MsgCountFreq = 10   // count word frequency in a file
const MsgCombineFreq = 11 // combine the word frequency from another collector

const MsgCollectorIdle = 50     // msg from collector to show that it is idle
const MsgCollectorDelivery = 51 // msg from collector to deliver data

const MsgDismissWorker = 100    // "kill" the collector
const MsgDeliverData = 101      // deliver the records from the collector
const MsgClearData = 102        // clear the records of the collector
const MsgSortAndSave2Disk = 103 // sort and save the records of the collector to disk
