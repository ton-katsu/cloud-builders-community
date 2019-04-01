package slackbot

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	cloudbuild "google.golang.org/api/cloudbuild/v1"
)

// Notify posts a notification to Slack that the build is complete.
func Notify(b *cloudbuild.Build, title string, icon string, webhook string) {
	url := fmt.Sprintf("https://console.cloud.google.com/cloud-build/builds/%s?project=%s", b.Id, b.ProjectId)
	var c string
	switch b.Status {
	case "SUCCESS":
		c = "#9CCC65"
	case "FAILURE", "CANCELLED":
		c = "#FF5252"
	case "STATUS_UNKNOWN", "INTERNAL_ERROR":
		c = "#FF5252"
	default:
		c = "#FF5252"
	}
	t := fmt.Sprintf("%s, Id: %s, Status: %s", title, b.Id, b.Status)
	j := fmt.Sprintf(
		`{	"icon_emoji": "%s",
			"username": "Cloud Build/%s",
			"attachments": [
				{
					"color": "%s",
					"title": "%s",
					"title_link": "%s",
				}
			]
		}`, icon, b.ProjectId, c, t, url)

	r := strings.NewReader(j)
	resp, err := http.Post(webhook, "application/json", r)
	if err != nil {
		log.Fatalf("Failed to post to Slack: %v", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("Posted message to Slack: [%v], got response [%s]", j, body)
}
