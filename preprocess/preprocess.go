package preprocess

import (
	"fmt"
	"os"
	"simple-implementation-mapreduce-invertedindex-v2/common"
	"strconv"
	"strings"
)

const splitSep = "CHAPTER "

// SplitOriginalTxt is used for splitting the original txt by chapters.
// Note that the index starts from 1.
func SplitOriginalTxt() {
	byt, err := os.ReadFile(common.PapTxtPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	s := string(byt)
	lists := strings.Split(s, splitSep)
	index := 1

	for _, str := range lists {
		if len(str) == 0 {
			continue
		}
		newFilename := common.PapDividedPathPrefix + strconv.Itoa(index) + common.TxtAppendix
		err = os.WriteFile(newFilename, []byte(str), 0644)
		if err != nil {
			fmt.Println("error for index ", index)
			fmt.Println(err)
		}
		index ++
	}
}
