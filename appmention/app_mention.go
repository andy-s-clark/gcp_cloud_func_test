package appmention

// When triggered by an "app_mention" event, publish the text to pub/sub

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/pubsub"
)

var topic *pubsub.Topic

func init() {
	projectId := os.Getenv("PROJECT_ID")
	if projectId == "" {{
		log.Fatalf("PROJECT_ID is required")
	}}

	topicName := os.Getenv("TOPIC_NAME")
	if topicName == "" {{
		log.Fatalf("TOPIC_NAME is required")
	}}

	ctx := context.Background()
	var client *pubsub.Client
	var err error
	client, err = pubsub.NewClient(ctx, projectId)
	if err != nil {
		log.Fatalf("pubsub.NewClient: %v", err)
	}

	topic = client.Topic(topicName)
	exists, err := topic.Exists(ctx)
	if err != nil {
		log.Fatalf("pubsub.NewClient.Topic: %v", err)
	}
	if !exists {
		log.Fatalf("Topic " + topicName + " not found")
	}
}

// Parsing the event in this function as a learning experience.
// Normally we'd use the github.com/slack-go/slack/slackevents package and verify tokens, etc.
type SlackAppMentionEvent struct {
	Type string `json:"type"`
	User string `json:"user"`
	Text string `json:"text"`
	Ts string `json:"ts"`
	Channel string `json:"channel"`
	EventTs string `json:"event_ts"`
}

type MentionMessage struct {
	Text string `json:"text"`
	Ts string `json:"ts"`
	Channel string `json:"channel"`
}

func AppMention(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "405 Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !strings.EqualFold(r.Header.Get("Content-Type"), "application/json") {
		http.Error(w, "406 Content type not acceptable", http.StatusNotAcceptable)
		return
	}
	var d SlackAppMentionEvent
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	if d.Type != "app_mention" {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	var mm = MentionMessage{ Text: d.Text, Ts: d.Ts, Channel: d.Channel }

	mmp, err := json.Marshal(mm)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	psm := &pubsub.Message {
		Data: mmp,
	}

	_, err = topic.Publish(r.Context(), psm).Get(r.Context())
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
