package main

import (
    "encoding/json"
    "fmt"
    "github.com/pelletier/go-toml"
    "gopkg.in/natefinch/lumberjack.v2"
    "log"
    "net/http"
    "os"
    "path"
)

// Represents entries in our log files.
type line struct {
    UserName string `json:"user_name"`
    Timestamp string `json:"timestamp"`
    Text string `json:"text"`
}

// Globals needed by the handler.
var secretToken string
var logDirectory string

// The web handler which receives messages and logs them.
func handler(w http.ResponseWriter, r *http.Request) {
    // Reject requests that don't have the shared secret with slack.
    if (r.Method != "POST" || r.FormValue("token") != secretToken) {
        w.WriteHeader(http.StatusForbidden)
        fmt.Fprintf(os.Stderr, "Bad request from %s\n", r.RemoteAddr, r)
        return
    }

    lineJson := &line{
        UserName: r.FormValue("user_name"),
        Timestamp: r.FormValue("timestamp"),
        Text: r.FormValue("text")}
    bytes, err := json.Marshal(lineJson)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error marshalling for request:\n", r)
        return
    }

    log.Println(string(bytes))
}


func main() {
    // Read config.toml
    config, err := toml.LoadFile("config.toml")
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error opening config file. Did you remember to `cp config.toml.example config.toml`?\n", err.Error())
        os.Exit(1)
    }
    secretToken = config.Get("slack-logger.secret_token").(string)
    logDirectory = config.Get("slack-logger.log_directory").(string)
    port := config.Get("slack-logger.port").(int64)

    // Set up our logger.
    log.SetOutput(&lumberjack.Logger{ Filename: path.Join(logDirectory, "slacklog.txt") })
    log.SetPrefix("")
    log.SetFlags(0)

    // Begin listening.
    http.HandleFunc("/", handler)
    err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error starting server: ", err)
        os.Exit(1)
    }
}
