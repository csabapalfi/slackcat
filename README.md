# slackcat

Post to Slack from your standard input.

## usage

```sh
echo "hello" | slackcat -token=$SLACK_TOKEN -channel=$YOUR_CHANNEL -tee
```

## installation

Grab the [latest release from Github](https://github.com/csabapalfi/slackcat/releases/latest).

## why?

There are lots of repos named slackcat (written in Go even).
The one I tried didn't work and it seemed too simple to implement hence the re-write.
