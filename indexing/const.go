package indexing

const msgCountFreq = 10  // count word frequency in a file
const msgCombineFreq = 11  // combine the word frequency from another collector

const msgCollectorIdle = 50  // msg from collector to show that it is idle
const msgCollectorBusy = 51 // msg from collector to show that it is busy
const msgCollectorDelivery = 52  // msg from collector to deliver data

const msgDismissWorker = 100  // "kill" the collector
const msgDeliverData = 101  // deliver the records from the collector
const msgClearData = 102  // clear the records of the collector
const msgSortAndSave2Disk = 103  // sort and save the records of the collector to disk
