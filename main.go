package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/slack-go/slack"
)

var (
	version = "HEAD"
	commit  = "UNKNOWN"
)

func GetEnv(key string, default_ string) string {
	value := os.Getenv(strings.ToLower(key))
	if value == "" {
		value = os.Getenv(strings.ToUpper(key))
	}
	if value == "" {
		value = default_
	}
	return value
}

func GetRequiredEnv(key string) string {
	value := GetEnv(key, "")
	if value == "" {
		fmt.Fprintf(os.Stderr, "Environment variable `%s` is required.\n", key)
		os.Exit(2)
	}
	return value
}

func GetMessage(aydURL *url.URL, target string) string {
	u, err := aydURL.Parse("status.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate status endpoint URL: %s\n", err)
		os.Exit(1)
	}

	resp, err := http.Get(u.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch message: %s\n", err)
		return ""
	}
	defer resp.Body.Close()

	var msg struct {
		Incidents []struct {
			Target  string `json:"target"`
			Message string `json:"message"`
		} `json:"current_incidents"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse status data: %s\n", err)
		return ""
	}

	for _, incident := range msg.Incidents {
		if incident.Target == target {
			return incident.Message
		}
	}
	fmt.Fprintf(os.Stderr, "No such incident information: %s\n", target)
	return ""
}

func main() {
	if len(os.Args) != 5 {
		fmt.Fprintln(os.Stderr, "$ ayd-slack-alert SLACK_URI TARGET_URI TARGET_STATUS TARGET_CHECKED_AT")
		os.Exit(2)
	}

	fmt.Printf("ayd-slack-alert %s (%s): ", version, commit)

	target := os.Args[2]
	status := os.Args[3]
	checkedAt, err := time.Parse(time.RFC3339, os.Args[4])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Environment variable `ayd_checked_at` is invalid: %s\n", err)
		os.Exit(2)
	}

	aydURL, err := url.Parse(GetEnv("ayd_url", "http://localhost:9000"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Environment variable `ayd_url` is invalid: %s\n", err)
		os.Exit(2)
	}
	statusPage, err := aydURL.Parse("status.html")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate status page URL: %s\n", err)
		os.Exit(1)
	}

	message := GetMessage(aydURL, target)

	attachmentStyle := "warning"
	if status == "FAILURE" {
		attachmentStyle = "danger"
	}

	webhookURL := GetRequiredEnv("slack_webhook_url")
	err = slack.PostWebhook(webhookURL, &slack.WebhookMessage{
		Attachments: []slack.Attachment{{
			Color:     attachmentStyle,
			Fallback:  message,
			Title:     fmt.Sprintf("[%s] %s", status, target),
			TitleLink: target,
			Text:      message,
			Footer:    aydURL.String(),
			Ts:        json.Number(strconv.FormatInt(checkedAt.Unix(), 10)),
			Actions: []slack.AttachmentAction{{
				Name: "Status Page",
				Text: "Status Page",
				Type: "button",
				URL:  statusPage.String(),
			}},
		}},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to send message: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Alert sent to Slack")
}
