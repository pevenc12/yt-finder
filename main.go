package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	// "time"

	"github.com/joho/godotenv"
	"github.com/pevenc12/yt-finder/helper"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

var period int

func main() {
	// parse parameters
	flag.Parse()
	terms, _, errParse := helper.ParseFlags(flag.Args())
	if errParse != nil {
		panic(errParse)
	}

	// parse environment variable
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	key := os.Getenv("YOUTUBE_API_KEY")

	// initialize a service
	client := &http.Client{
		Transport: &transport.APIKey{Key: key},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	// query videos
	call := service.Search.List([]string{"id", "snippet"}).
		// PublishedAfter(time.Now().Add(day * 24 * time.Hour()).Format()).
		// Q(`coding | remote`).
		Q(strings.Join(terms, "|")).
		MaxResults(50)
	response, err := call.Do()
	handleError(err, "")

	videos := make(map[string]string)
	channels := make(map[string]string)
	playlists := make(map[string]string)

	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			videos[item.Id.VideoId] = item.Snippet.Description
		case "youtube#channel":
			channels[item.Id.ChannelId] = item.Snippet.Title
		case "youtube#playlist":
			playlists[item.Id.PlaylistId] = item.Snippet.Title
		}
	}

	printIDs("Videos", videos)
	printIDs("Channels", channels)
	printIDs("Playlists", playlists)

}

func printIDs(sectionName string, matches map[string]string) {
	fmt.Printf("%v:\n", sectionName)
	for id, title := range matches {
		fmt.Printf("[%v] %v\n", id, title)
	}
	fmt.Printf("\n\n")
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}
