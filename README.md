# slackcat

Post to Slack from stdin.

## usage examples

* simple usage
```sh
echo "hello" | slackcat -channel=$YOUR_CHANNEL
```

* use token from environment variable
```sh
export SLACK_TOKEN=$SLACK_TOKEN
echo "hello" | slackcat -channel=$YOUR_CHANNEL
```

* use token from command line flag
```sh
echo "hello" | slackcat -channel=$YOUR_CHANNEL -token=$SLACK_TOKEN
```

* treat stdin as [Slack attachment](https://api.slack.com/docs/message-attachments)
```sh
echo attachment.json | slackcat -channel=$YOUR_CHANNEL --attachment
```
Check out the [Slack Message Builder](https://api.slack.com/docs/messages/builder) to give you an idea. Also please find a simple example JSON below:
```json
{
  "fallback": "Required plain-text summary of the attachment.",
  "color": "#36a64f",
  "title": "Slack API Documentation",
  "title_link": "https://api.slack.com/",
  "text": "Optional text that appears within the attachment",
}
```

* tee stdin to both stdout and Slack
```sh
echo "hello" | slackcat -channel=$YOUR_CHANNEL -tee
```

* verbose output (print slack API response)
```sh
echo "hello" | slackcat -channel=$YOUR_CHANNEL -v
```

## installation

Grab the [latest release from Github](https://github.com/csabapalfi/slackcat/releases/latest).

## why?

There are lots of repos named slackcat (written in Go even).
The one I tried didn't work and it seemed too simple to implement hence the re-write.
