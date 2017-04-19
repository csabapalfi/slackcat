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

func buildSlackURL(channel, message, token string) string {
	slackURL := &url.URL{
		Host:   "slack.com",
		Scheme: "https",
		Path:   "api/chat.postMessage",
	}
	q := slackURL.Query()
	q.Set("token", token)
	q.Set("channel", channel)
	q.Set("text", message)
	q.Set("as_user", "true")
	slackURL.RawQuery = q.Encode()
	return slackURL.String()
}

func postToSlack(channel, message, token string, verbose bool) {
	type SlackResult struct {
		Ok    bool
		Error string
	}
	url := buildSlackURL(channel, message, token)
	result := &SlackResult{}
	err := getJSON(url, result, verbose)
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
	token := flag.String("token", "", "Slack token")
	channel := flag.String("channel", "", "Slack channel")
	tee := flag.Bool("tee", false, "tee stdin to both stdout and slack")

	flag.Parse()

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

	postToSlack(*channel, message, *token, *verbose)
}
