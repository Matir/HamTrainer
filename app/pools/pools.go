package pools

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

var classes = []string{"technician", "general", "extra"}

// Data about a single question
type Question struct {
	Topic string
	Number string
	Citation string
	Question string
	Correct string
	Answers map[string]string
}

// Data about a pool
type QuestionPool struct {
	version int
	Subelements map[string]string
	Topics map[string]string
	Questions map[string]Question
}

// Cache data once it is loaded once
var questionPools map[string]*QuestionPool = make(map[string]*QuestionPool)

// Base path for json files
var jsonBase = "data/json"

// Get data for the question pool specified by Pool
func GetPool(class string) (*QuestionPool, error) {
	if val, ok := questionPools[class]; ok {
		return val, nil
	}
	return loadPool(class)
}

// Set the directory to look for JSON files in
func SetJsonBase(path string) {
	jsonBase = path
}

// Load a single pool from the relevant file
func loadPool(class string) (*QuestionPool, error) {
	if !isValidClass(class) {
		return nil, fmt.Errorf("Invalid class specified: %s", class)
	}

	poolText, err := getPoolSource(class)
	if err != nil {
		return nil, fmt.Errorf("Error loading pool source:", err)
	}

	pool := new(QuestionPool)
	// TODO implement versioned JSON loading
	err = json.Unmarshal(poolText, pool)
	if err != nil {
		return nil, fmt.Errorf("Error parsing pool JSON:", err)
	}

	questionPools[class] = pool
	return pool, nil
}

// Load a pool source file
func getPoolSource(class string) ([]byte, error) {
	fname := filepath.Join(jsonBase, class + ".json")
	return ioutil.ReadFile(fname)
}

// Check to see if the class is a valid class
func isValidClass(class string) bool {
	for _, val := range classes {
		if val == class {
			return true
		}
	}
	return false
}
