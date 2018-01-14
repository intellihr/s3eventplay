package s3eventplay

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func listEventFileKeys(svc *s3.S3, bucket string, prefix string) ([]string, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}

	var contents []*s3.Object
	err := svc.ListObjectsV2Pages(input,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			contents = append(contents, page.Contents...)
			return !lastPage
		})
	if err != nil {
		return nil, err
	}

	keys := make([]string, len(contents))
	for i, content := range contents {
		keys[i] = *content.Key
	}
	sort.Strings(keys)
	return keys, nil
}

func listEventFileKeysForParams(svc *s3.S3, params Params) ([]string, error) {
	prefixes, err := params.s3Prefixes()
	if err != nil {
		return prefixes, err
	}
	allKeys := []string{}
	for _, prefix := range prefixes {
		keys, err := listEventFileKeys(svc, params.Bucket, prefix)
		if err != nil {
			return keys, err
		}
		allKeys = append(allKeys, keys...)
	}

	return allKeys, nil
}

func getEventFile(svc *s3.S3, bucket string, key string) (string, error) {
	result, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", err
	}
	defer result.Body.Close()

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, result.Body)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func processEventsString(content string) (err error) {
	events := strings.Split(content, "\n")
	for _, event := range events {
		if event == "" {
			continue
		}
		data, err := base64.StdEncoding.DecodeString(event)
		if err != nil {
			err = fmt.Errorf("error decoding: %s", event)
			break
		}
		fmt.Printf("%s\n", data)
	}
	return err
}

type Params struct {
	Bucket    string
	Stream    string
	BatchSize int
	Dates     string
}

func (params *Params) s3Prefixes() ([]string, error) {
	dates, err := params.dateList()
	if err != nil {
		return dates, err
	}
	results := make([]string, len(dates))
	for i, date := range dates {
		results[i] = fmt.Sprintf("%s/%s", params.Stream, date)
	}
	return results, nil
}

func (params *Params) dateList() ([]string, error) {
	dates := strings.Split(params.Dates, "~")
	if len(dates) <= 1 {
		return []string{dates[0]}, nil
	}
	start, err := time.Parse("2006-01-02", dates[0])
	if err != nil {
		return []string{}, err
	}
	end, err := time.Parse("2006-01-02", dates[1])
	if err != nil {
		return []string{}, err
	}
	end = end.AddDate(0, 0, 1)
	results := []string{}

	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		results = append(results, d.Format("2006-01-02"))
	}
	return results, nil
}

func PlayEvents(params Params) {
	svc := s3.New(session.New())

	keys, err := listEventFileKeysForParams(svc, params)
	if err != nil {
		panic(err)
	}
	inputs := make([]interface{}, len(keys))
	for i, key := range keys {
		inputs[i] = key
	}

	SequentialBatch(inputs,
		func(input interface{}) interface{} {
			content, err := getEventFile(svc, params.Bucket, input.(string))
			if err != nil {
				fmt.Fprintln(os.Stderr, "error downloading "+input.(string))
				panic(err)
			}
			return content
		}, func(result interface{}) {
			err := processEventsString(result.(string))
			if err != nil {
				panic(err)
			}
		}, params.BatchSize)
}
