package slackbot

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	cloudbuild "google.golang.org/api/cloudbuild/v1"
)

// Notify posts a notification to Slack that the build is complete.
func Notify(b *cloudbuild.Build, title string, icon string, tag string, webhook string) {
	burl := fmt.Sprintf("https://console.cloud.google.com/cloud-build/builds/%s?project=%s", b.Id, b.ProjectId)
	query := fmt.Sprintf("tags=\"%s\"", tag)
	params := url.Values{}
	params.Add("query", query)
	params.Add("project", b.ProjectId)
	turl := fmt.Sprintf("https://console.cloud.google.com/cloud-build/builds?%s", params.Encode())
	text := fmt.Sprintf("[%s] %s. Id: %s\nStartTime: %s\nFinishTime: %s", b.Status, title, b.Id, b.StartTime, b.FinishTime)
	var c string
	switch b.Status {
	case "SUCCESS":
		c = "#9CCC65"
	case "FAILURE":
		c = "#FF5252"
	case "CANCELLED":
		c = "#CCD1D9"
	case "STATUS_UNKNOWN", "INTERNAL_ERROR":
		c = "#FF5252"
	default:
		c = "#FF5252"
	}
	j := fmt.Sprintf(
		`{	"icon_emoji": "%s",
			"username": "Cloud Build/%s",
			"attachments": [
				{
					"color": "%s",
					"text": "%s",
					"actions": [
						{
							"type": "button",
							"text": "Details",
							"url": "%s"
						},
						{
							"type": "button",
							"text": "Results using %s tag",
							"url": "%s"
						}
					]
				}
			]
		}`, icon, b.ProjectId, c, text, burl, tag, turl)

	r := strings.NewReader(j)
	resp, err := http.Post(webhook, "application/json", r)
	if err != nil {
		log.Fatalf("Failed to post to Slack: %v", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("Posted message to Slack: [%v], got response [%s]", j, body)
}
