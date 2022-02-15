package search

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"simple-implementation-mapreduce-invertedindex-v2/common"
	"simple-implementation-mapreduce-invertedindex-v2/refine"
	"sort"
	"strconv"
	"strings"
)

const errInput = "invalid input not recognized"

type Engine struct {
	RecordsPath 	string
	Records 		common.Records
	AllWords 		[]string  			// a list of all the words, sorted
	Snapshots		refine.Snapshots
	online 			bool
}

// NewSearchEngine returns a pointer to a new Engine object. The Engine object is already initialized.
func NewSearchEngine() (*Engine, error) {
	engine := &Engine{
		RecordsPath: common.InvertedIndexSavePath,
		Records: common.NewRecords(),
		Snapshots: make(refine.Snapshots),
		online: true,
	}

	// load the inverted indexes
	bytes, err := os.ReadFile(common.InvertedIndexSavePath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &engine.Records)
	if err != nil {
		return nil, err
	}

	// count all the words
	allWords := make([]string, len(engine.Records))
	index := 0
	for word, _ := range engine.Records {
		allWords[index] = word
		index ++
	}
	sort.Strings(allWords)
	engine.AllWords = allWords

	// load the snapshots
	bytes, err = os.ReadFile(common.SnapshotsSavePath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &engine.Snapshots)
	if err != nil {
		return nil, err
	}

	return engine, nil
}

// Run executes the main loop.
func (e *Engine) Run() {
	scanner := bufio.NewScanner(os.Stdin)
	for e.online {
		e.welcomeText()

		scanner.Scan()
		userInput := scanner.Text()
		good := false

		switch userInput {
		case "count":
			good = true
			e.output(len(e.AllWords))

		case "all":
			good = true
			e.output(e.AllWords)

		default:
			args := strings.Split(userInput, " ")
			if len(args) > 1 {
				switch args[0] {
				case "search":
					good = true
					searched := e.search(args[1:])
					e.output(searched)

				case "vague":
					e.output("not implemented")

				case "command":
					if args[1] == "exit" {
						good = true
						e.online = false
					}
				}
			}
		}

		if !good {
			e.output(errInput)
		}
	}
	return
}

// prints the welcome texts on the screen.
func (e *Engine) welcomeText() {
	fmt.Println("Please choose your next action.")
	fmt.Println("type 'count' to view the number of words in the database.")
	fmt.Println("Type 'all' to view all the words.")
	fmt.Println("Type 'search [word]...' to view the statistics of some certain words.")
	fmt.Println("Type 'vague [word]' to find likely words. NOT IMPLEMENTED YET")
	fmt.Println("type 'command exit' to exit the program.")
}

// outputs the text onto the user interface; currently it is merely a wrap of fmt.Println().
func (e *Engine) output(a interface{}) {
	switch a.(type) {
	case []*singleFileSummary:
		sfss := a.([]*singleFileSummary)
		conclusion := "a total of " + strconv.Itoa(len(sfss)) + " results found.\n"
		fmt.Println(conclusion)
		for _, sfs := range sfss {
			fmt.Println(sfs)
			fmt.Println(e.Snapshots[sfs.filename] + "...\n")
		}
		fmt.Println(conclusion)

	default:
		fmt.Println(a)
	}
}

// searches for a list of words and return the relevant results.
func (e *Engine) search(words []string) interface{} {
	summary := NewSummary()
	for _, word := range words {
		for _, record := range e.Records[word] {
			summary.Add(word, record)
		}
	}
	return summary.Sort()
}
