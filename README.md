# Slack Logger
Slack logger is a simple webserver that listens for messages from an outgoing Slack webhook and logs them. This can be useful for saving chat history
from a free slack instance.

The log format is one json object per line, each containing a string `user_name`, a float `timestamp` and a string `message`. Example:
 
    {"user_name":"blackpatch","timestamp":"1441578779.017708","text":"How To Dismantle An Atomic Bomb, like in that U2 album, Achtung Baby"}

The logs are rotated every 100 MB using `lumberjack` (https://github.com/natefinch/lumberjack).

## Installation
 * Crreate an outgoing webhook integration in slack and note the secret token.
 * Copy `config.toml.example` to `config.toml`.
 * Fill in the fields of `config.toml` with your secret token, the location where you want to save your logs, and the port you want slack-logger to listen on.
 * `go run slack-logger.go`
 * In the outgoing webhook integration, fill in the url for the machine for your slack-logger. Remember to include the port.

## Known Issues
 * By design slack-logger only logs user message events.
 * slack-logger uses usernames rather than ids, so name changes would need to be tracked separately if real identity is important.
 * The connection is not SSL yet. http://golang.org/src/crypto/tls/generate_cert.go could be used for this.
 * There should be more logging around attempted connections.