package main

import (
	"bytes"
	"encoding/json"
	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type subreddit struct {
	Name  string   `json:"name"`
	Words []string `json:"words"`
}

type announcer struct {
	SubNSearches map[string][]string
	Slack        string
}

//  graw.Scan() will use this Post handler
func (a *announcer) Post(post *reddit.Post) error {

	//Grab list of search words by subreddit
	filter := a.SubNSearches[post.Subreddit]

	//Iterate over list looking for match and then post to slack if match true
	for _, word := range filter {
		if strings.Contains(post.Title, word) {
			msg := "This post looks interesting, check it out: " + post.Title + string(" \n") + post.URL
			postMessage(msg, a.Slack)
			return nil
		}
	}

	return nil
}

func main() {

	// Get an api handle to reddit for a logged out (script) program,
	// which forwards this user agent on all requests and issues a request at
	// most every 5 seconds.
	apiHandle, err := reddit.NewScript("Ubuntu:github.com/AnthonyLaiuppa/snooper:v0.0.1 (by /u/ThisIsMyRedditAccount)", 5*time.Second)
	if err != nil {
		log.Fatalln("Failed to create NewScript: ", err)
		return
	}

	//Read in what subs we want to monitor with what keywords
	a, err := setUp()
	if err != nil {
		log.Fatalln("Failed to map Subreddits to keywords: ", err)
		return
	}

	//Grab all the map keys to use as our subreddits to open streams to
	keys := make([]string, 0, len(a.SubNSearches))
	for k := range a.SubNSearches {
		keys = append(keys, k)
	}

	// Create a configuration specifying what event sources on Reddit graw
	// should connect to the bot.
	cfg := graw.Config{Subreddits: keys}

	// launch a graw scan in a goroutine using the bot, handle, and config. The
	// returned "stop" and "wait" are functions. "stop" will stop the graw run
	// at any time, and "wait" will block until it finishes.
	stop, wait, err := graw.Scan(&a, apiHandle, cfg)

	// This time, let's block so the bot will announce (ideally) forever.
	if err := wait(); err != nil {
		log.Fatalln("graw run encountered an error: %v \n", err)
		stop()
	}
}

func postMessage(msg string, url string) {

	//Marshal up some json to post
	message := map[string]interface{}{"text": msg}
	mib, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}

	//Make post request to slack webhook
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(mib))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	return
}

func setUp() (a announcer, err error) {

	//Get slack webhook URL
	a.Slack = os.Getenv("SWHURL")

	//Open the json file containing our search parameters
	jsonFile, err := os.Open("./sns.ini")
	if err != nil {
		log.Fatalln("File not found: ", err)
		return a, err
	}
	defer jsonFile.Close()

	var subreddits []subreddit

	//Read in our JSON
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalln("File unreadable: ", err)
		return a, err
	}

	//Attempt to unmarshal our json
	err = json.Unmarshal(byteValue, &subreddits)
	if err != nil {
		log.Fatalln("Error unmarshaling JSON : ", err)
		return a, err
	}

	//Create a map sub->keywords
	a.SubNSearches = make(map[string][]string)
	for _, sub := range subreddits {
		a.SubNSearches[sub.Name] = sub.Words
	}
	return a, err
}
