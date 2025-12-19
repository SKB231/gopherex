package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/SKB231/gopherex/ex13/hn"
)

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()
	fmt.Println("Building template index.gohtml")
	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl))

	fmt.Println("Starting listener..")
	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

type cacheStruct struct {
	vals       []int
	recordedAt time.Time
}

var cache cacheStruct

const THRESHOLD_DURATION time.Duration = time.Second * 10

var topItemsMutex sync.Mutex

func retriveTopItems(client hn.Client) (ids []int, err error) {

	topItemsMutex.Lock()
	defer topItemsMutex.Unlock()
	if len(cache.vals) == 0 || time.Since(cache.recordedAt) >= THRESHOLD_DURATION {
		//fmt.Println(ok, &client, cache, time.Since(cache.recordedAt), time.Since(resp.recordedAt) >= THRESHOLD_DURATION)
		ids, err = client.TopItems()
		fmt.Println("USED CLIENT!")
		if err != nil {
			fmt.Println("Error when retrieving the top ids..", err)
			return
		}
		cache = struct {
			vals       []int
			recordedAt time.Time
		}{
			vals:       ids,
			recordedAt: time.Now(),
		}
		fmt.Println("Recorded to cache..")
		return
	}

	fmt.Println("USED CACHE!")
	return cache.vals, nil
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var client hn.Client
		ids, err := retriveTopItems(client)
		if err != nil {
			http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
			return
		}
		type rankedItem struct {
			rank int
			item
		}
		var stories []item

		storyBuffer := make(chan rankedItem, numStories)
		//endSearch := make(chan bool) // Done with finding numStory stories..
		searchEnded := false
		var wg sync.WaitGroup
		for rank, id := range ids {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				if searchEnded {
					return
				}
				hnItem, err := client.GetItem(id)
				if err != nil {
					return
				}
				if searchEnded {
					return
				}

				item := parseHNItem(hnItem)
				if searchEnded {
					return
				}
				if isStoryLink(item) {
					//stories = append(stories, item)
					storyBuffer <- rankedItem{
						rank: rank,
						item: item,
					}
					if len(storyBuffer) >= numStories {
						//
						searchEnded = true
					}
				}
			}(id)
		}
		wg.Wait()
		rankedStories := make([]rankedItem, 0)

		for {
			if len(storyBuffer) <= 0 {
				break
			}
			nextStory := <-storyBuffer
			rankedStories = append(rankedStories, nextStory)
		}

		sort.Slice(rankedStories, func(i, j int) bool {
			return rankedStories[i].rank < rankedStories[j].rank
		})

		for _, story := range rankedStories {
			stories = append(stories, story.item)
		}

		data := templateData{
			Stories: stories,
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}
