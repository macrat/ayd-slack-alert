package main

import (
	"encoding/json"
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
		logger.Failure(fmt.Sprintf("Environment variable `%s` is required", key))
		os.Exit(0)
	}
	return value
}

func GetMessage(aydURL, targetURL *url.URL) (string, error) {
	resp, err := ayd.Fetch(aydURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch status: %w", err)
	}

	rs, err := resp.RecordsOf(targetURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch status: %w", err)
	}

	return rs[len(rs)-1].Message, nil
}

func main() {
	args, err := ayd.ParseAlertPluginArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, "$ ayd-slack-alert MAILTO_URI TARGET_URI TARGET_STATUS TARGET_CHECKED_AT")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	logger := ayd.NewLogger(args.AlertURL)

	webhookURL := GetRequiredEnv(logger, "slack_webhook_url")

	aydURL, err := url.Parse(GetEnv("ayd_url", "http://localhost:9000"))
	if err != nil {
		logger.Failure(fmt.Sprintf("environment variable `ayd_url` is invalid: %s", err))
		return
	}
	statusPage, err := aydURL.Parse("status.html")
	if err != nil {
		logger.Failure(fmt.Sprintf("failed to generate status page URL: %s", err))
		return
	}

	message, err := GetMessage(aydURL, args.TargetURL)
	if err != nil {
		logger.Unknown(err.Error())
	}

	attachmentStyle := "warning"
	if args.Status == ayd.StatusFailure {
		attachmentStyle = "danger"
	}

	err = slack.PostWebhook(webhookURL, &slack.WebhookMessage{
		Attachments: []slack.Attachment{{
			Color:     attachmentStyle,
			Fallback:  message,
			Title:     fmt.Sprintf("[%s] %s", args.Status.String(), args.TargetURL.String()),
			TitleLink: args.TargetURL.String(),
			Text:      message,
			Footer:    aydURL.String(),
			Ts:        json.Number(strconv.FormatInt(args.CheckedAt.Unix(), 10)),
			Actions: []slack.AttachmentAction{{
				Name: "Status Page",
				Text: "Status Page",
				Type: "button",
				URL:  statusPage.String(),
			}},
		}},
	})
	if err != nil {
		logger.Failure(fmt.Sprintf("failed to send message: %s", err))
		return
	}
	logger.Healthy("Alert sent to Slack")
}
