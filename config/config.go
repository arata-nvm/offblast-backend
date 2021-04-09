package config

import "os"

func Port() string {
	return os.Getenv("PORT")
}

func SentryDsn() string {
	return os.Getenv("SENTRY_DSN")
}

func SlackWebhook() string {
	return os.Getenv("SLACK_WEBHOOK")
}
