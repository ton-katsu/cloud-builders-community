// Post build status results to Slack.

package main

import (
	"context"
	"flag"
	"log"

	"./slackbot"
)

var (
	title   = flag.String("title", "", "Notification title. e.g. trigger name.")
	icon    = flag.String("icon", "", "Notification icon.")
	build   = flag.String("build", "", "Build ID being monitored")
	webhook = flag.String("webhook", "", "Slack webhook URL")
	mode    = flag.String("mode", "trigger", "Mode the builder runs in")
)

func main() {
	log.Print("Starting slackbot")
	flag.Parse()
	ctx := context.Background()

	if *title == "" {
		log.Fatalf("title must be provided.")
	}

	if *icon == "" {
		log.Fatalf("icon must be provided.")
	}

	if *webhook == "" {
		log.Fatalf("Slack webhook must be provided.")
	}

	if *build == "" {
		log.Fatalf("Build ID must be provided.")
	}

	if *mode == "trigger" {
		// Trigger another build to run the monitor.
		log.Printf("Starting trigger mode for build %s", *build)
		slackbot.Trigger(ctx, *title, *icon, *build, *webhook)
		return
	}
	if *mode == "monitor" {
		// Monitor the other build until completion.
		log.Printf("Starting monitor mode for build %s", *build)
		slackbot.Monitor(ctx, *title, *icon, *build, *webhook)
		return
	}
	log.Fatalf("Mode must be provided.")
}
