package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type ArcUnit struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type ArcMap map[string]ArcUnit

func (ac ArcMap) ArcHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("layout.html"))
		path := r.RequestURI
		if path == "/" {
			path = "/intro"
		}
		path = path[1:]
		if value, ok := ac[path]; ok {
			tmpl.Execute(w, value)
			return
		} else {
			http.NotFound(w, r)
			return
		}
	}
}

func main() {
	testStory := processStoryFromJSON()
	http.ListenAndServe(":8080", testStory.ArcHandler())
}

func processStoryFromJSON() *ArcMap {
	// Open the JSON file
	file, err := os.Open("story.json")
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	// Read the file's content into a byte slice
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	// Unmarshal the JSON data into the struct
	var arcStory ArcMap
	err = json.Unmarshal(bytes, &arcStory)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %s", err)
	}
	//fmt.Println(arcStory)
	return &arcStory
}
