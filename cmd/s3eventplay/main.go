package main

import (
	"os"

	"github.com/urfave/cli"
	"github.com/intellihr/s3eventplay"
)

func main() {
	app := cli.NewApp()
	app.Name = "s3eventplay"
	app.Version = "master"
	app.Author = "Soloman Weng (soloman1124@gmail.com)"
	app.Usage = "This tool plays json events from s3 bucket"
	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "bucket, b",
			Usage: "target AWS S3 bucket, required",
			EnvVar: "S3_BUCKET",
		},
		cli.StringFlag{
			Name: "stream, s",
			Usage: "target kinesis stream name or folder name that stores the event files, required",
			EnvVar: "EVENT_STREAM",
		},
		cli.IntFlag{
			Name: "batch, n",
			Value: 5,
			Usage: "target batch size (i.e. number of concurrent file downloads).",
			EnvVar: "BATCH_SIZE",
		},
		cli.StringFlag{
			Name: "dates, d",
			Usage: "target dates (date range) to play the events (e.g. YYYY-MM-DD or YYYY-MM-DD~YYYY-MM-DD)",
			EnvVar: "EVENT_DATES",
		},
	}
	app.Action = func(c *cli.Context) error {
		bucket := c.String("bucket")
		if bucket == "" {
			return cli.NewExitError("[bucket] is required", 1)
		}

		stream := c.String("stream")
		if stream == "" {
			return cli.NewExitError("[stream] is required", 1)
		}

		dates := c.String("dates")
		if dates == "" {
			return cli.NewExitError("[dates] is required", 1)
		}

		batchSize := c.Int("batch")

		s3eventplay.PlayEvents(s3eventplay.Params{
			Bucket: bucket,
			Stream: stream,
			BatchSize: batchSize,
			Dates: dates,
		})
		return nil
	}
	app.Run(os.Args)
}
