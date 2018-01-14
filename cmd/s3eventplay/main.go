package main

import (
	"github.com/intellihr/s3eventplay"
)

func main() {
	s3eventplay.PlayEvents(s3eventplay.Params{
		Bucket: "event-bus-service-soloman-event-backup",
		Stream: "event-bus-service-soloman-stream",
		BatchSize: 10,
	})
}
