package search

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"simple-implementation-mapreduce-invertedindex-v2/common"
	"sort"
	"strings"
)

const errInput = "invalid input not recognized"

type Engine struct {
	RecordsPath 	string
	Records 		common.Records
	AllWords 		[]string  			// a list of all the words, sorted
	online 			bool
}

// NewSearchEngine returns a pointer to a new Engine object. The Engine object is already initialized.
//
// recordsPath: the path of the records (json file).
func NewSearchEngine(recordsPath string) (*Engine, error) {
	engine := &Engine{RecordsPath: recordsPath, Records: common.NewRecords(), online: true}
	bytes, err := os.ReadFile(recordsPath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &engine.Records)
	if err != nil {
		return nil, err
	}

	allWords := make([]string, len(engine.Records))
	index := 0
	for word, _ := range engine.Records {
		allWords[index] = word
		index ++
	}
	sort.Strings(allWords)
	engine.AllWords = allWords

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
	fmt.Println("Type 'search [word]' to view the statistics of a certain word.")
	fmt.Println("type 'command exit' to exit the program.")
}

// outputs the text onto the user interface; currently it is merely a wrap of fmt.Println().
func (e *Engine) output(a interface{}) {
	fmt.Println(a)
}

// searches for a list of words and return the relevant results.
func (e *Engine) search(words []string) interface{} {
	if len(words) == 1 {
		return e.Records[words[0]]
	} else {
		// todo
		return "not implemented"
	}
}
