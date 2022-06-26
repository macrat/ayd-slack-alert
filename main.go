package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/macrat/ayd/lib-ayd"
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

func GetRequiredEnv(logger ayd.Logger, key string) string {
	value := GetEnv(key, "")
	if value == "" {
		logger.Failure(fmt.Sprintf("Environment variable `%s` is required", key), nil)
		os.Exit(0)
	}
	return value
}

func Usage() {
	fmt.Fprintln(os.Stderr, "Usage: ayd-slack-alert SLACK_ALERT_URL CHECKED_AT TARGET_STATUS LATENCY TARGET_URL MESSAGE EXTRA_VALUES")
}

func main() {
	showVersion := flag.Bool("v", false, "show version and exit.")
	flag.Usage = Usage
	flag.Parse()

	if *showVersion {
		fmt.Printf("ayd-slack-alert %s (%s)\n", version, commit)
		return
	}

	args, err := ayd.ParseAlertPluginArgs()
	if err != nil {
		Usage()
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	logger := ayd.NewLogger(&ayd.URL{Scheme: args.AlertURL.Scheme})

	webhookURL := GetRequiredEnv(logger, "slack_webhook_url")

	aydURL, err := url.Parse(GetEnv("ayd_url", "http://localhost:9000"))
	if err != nil {
		logger.Failure(fmt.Sprintf("environment variable `ayd_url` is invalid: %s", err), nil)
		return
	}
	statusPage, err := aydURL.Parse("status.html")
	if err != nil {
		logger.Failure(fmt.Sprintf("failed to generate status page URL: %s", err), nil)
		return
	}

	var attachmentStyle string
	switch args.Status {
	case ayd.StatusHealthy:
		attachmentStyle = "good"
	case ayd.StatusFailure:
		attachmentStyle = "danger"
	default:
		attachmentStyle = "warning"
	}

	status := args.Status.String()
	if args.Status == ayd.StatusHealthy {
		status = "RESOLVED"
	}

	err = slack.PostWebhook(webhookURL, &slack.WebhookMessage{
		Attachments: []slack.Attachment{{
			Color:     attachmentStyle,
			Fallback:  args.Message,
			Title:     fmt.Sprintf("[%s] %s", status, args.TargetURL.String()),
			TitleLink: args.TargetURL.String(),
			Text:      args.Message,
			Footer:    aydURL.String(),
			Ts:        json.Number(strconv.FormatInt(args.Time.Unix(), 10)),
			Actions: []slack.AttachmentAction{{
				Name: "Status Page",
				Text: "Status Page",
				Type: "button",
				URL:  statusPage.String(),
			}},
		}},
	})
	if err != nil {
		logger.Failure(fmt.Sprintf("failed to send message: %s", err), nil)
		return
	}
	logger.Healthy("Alert sent to Slack", nil)
}
