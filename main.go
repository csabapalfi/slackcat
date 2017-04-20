package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

func debug(data []byte, err error) {
	if err == nil {
		fmt.Fprintf(os.Stderr, "%s\n\n", data)
	} else {
		fmt.Fprintf(os.Stderr, "%s\n\n", err)
	}
}

func getJSON(url string, target interface{}, verbose bool) error {
	response, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if verbose {
		debug((httputil.DumpResponse(response, true)))
	}
	return json.NewDecoder(response.Body).Decode(target)
}

func buildSlackURL(channel, message, token string, attachment bool) string {
	slackURL := &url.URL{
		Host:   "slack.com",
		Scheme: "https",
		Path:   "api/chat.postMessage",
	}
	q := slackURL.Query()
	q.Set("token", token)
	q.Set("channel", channel)
	q.Set("as_user", "true")
	if attachment {
		q.Set("attachments", "["+message+"]")
	} else {
		q.Set("text", message)
	}
	slackURL.RawQuery = q.Encode()
	return slackURL.String()
}

func postToSlack(slackURL string, verbose bool) {
	type SlackResult struct {
		Ok    bool
		Error string
	}
	result := &SlackResult{}
	err := getJSON(slackURL, result, verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n\n", err)
		os.Exit(1)
	}
	if err == nil && !result.Ok {
		fmt.Fprintf(os.Stderr, "%s\n\n", result.Error)
		os.Exit(1)
	}
}

func main() {
	verbose := flag.Bool("v", false, "verbose output")
	attachment := flag.Bool("attachment", false, "treat input as Slack attachment")
	token := flag.String("token", "", "Slack token")
	channel := flag.String("channel", "", "Slack channel")
	tee := flag.Bool("tee", false, "tee stdin to both stdout and Slack")

	flag.Parse()

	tokenFromEnv := os.Getenv("SLACK_TOKEN")

	var stdin io.Reader = os.Stdin
	if *tee {
		stdin = io.TeeReader(os.Stdin, os.Stdout)
	}
	var buffer bytes.Buffer
	scanner := bufio.NewScanner(stdin)
	for scanner.Scan() {
		buffer.WriteString(scanner.Text())
	}
	message := buffer.String()

	if *token == "" {
		token = &tokenFromEnv
	}

	slackURL := buildSlackURL(*channel, message, *token, *attachment)
	postToSlack(slackURL, *verbose)
}
