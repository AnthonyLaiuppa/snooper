package main

import (
	"testing"
	"encoding/json"
	"strings"
	"gotest.tools/assert"
)

type subredditTest struct {
	Name  string   `json:"name"`
	Words []string `json:"words"`
}

var exampleJson = []byte(`[{"name":"funny","words":["The","My","the","my","A","i","a","I"]},{"name":"pics","words":["beautiful","the","The","my","My","owl","cat","i","A","I"]}]`) 


//Validating our example data/structure with an unmarshal
func TestSetup(t *testing.T) {
	
	var subreddits []subredditTest

	err := json.Unmarshal(exampleJson, &subreddits)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON %v: ", err)
	}

}

//Tests our "imported data" for presence of search string
func TestPost(t *testing.T) {


	var subreddits []subredditTest

	err := json.Unmarshal(exampleJson, &subreddits)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON : %v", err)
	}

	//Create our map
	var a = make(map[string][]string)
	for _, sub := range subreddits {
		a[sub.Name] = sub.Words
	}

	//Validate our control word is present
	words := strings.Join(a["pics"], " ")
	assert.Assert(t, strings.Contains(words, "beautiful") == true)

}