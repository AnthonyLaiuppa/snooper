package main

import (
	"encoding/json"
	"github.com/turnage/graw/reddit"
	"gotest.tools/assert"
	"os"
	"testing"
)

var exampleJson = []byte(`[{"name":"funny","words":["The","My","the","my","A","i","a","I"]},{"name":"pics","words":["beautiful","the","The","my","My","owl","cat","i","A","I"]}]`)

//Validating our example data/structure with an unmarshal
func TestJson(t *testing.T) {

	var subreddits []subreddit

	err := json.Unmarshal(exampleJson, &subreddits)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON %v: ", err)
	}

}

// Call the post function and look for errors...
func TestPost(t *testing.T) {

	var post = reddit.Post{
		Title: "Testing is beautiful",
		URL:   "https://reddit.com/r/golang",
	}

	//Mock the struct that wouldve come from reading in our jsonfile
	var subreddits []subreddit
	err := json.Unmarshal(exampleJson, &subreddits)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON : %v", err)
	}

	//Instantiate an announcer struct
	var a announcer
	a.Slack = os.Getenv("SWHURL")
	a.SubNSearches = make(map[string][]string)
	for _, sub := range subreddits {
		a.SubNSearches[sub.Name] = sub.Words
	}

	err = a.Post(&post)
	assert.Assert(t, err == nil)
}

//Call postMessage and validate return code is 200
func TestPostMessage(t *testing.T) {
	swhurl := os.Getenv("SWHURL")
	msg := "This post looks like a test, check it out: https://reddit.com/r/golang"
	s := postMessage(msg, swhurl)
	assert.Assert(t, s == 200)
}

// Call setup and look for errors...
func TestSetup(t *testing.T) {
	_, err := setUp()
	if err != nil {
		t.Fatalf("Setup test failed: %v ", err)
		return
	}
	assert.Assert(t, err == nil)

}
