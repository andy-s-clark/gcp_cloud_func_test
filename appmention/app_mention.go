package appmention

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strings"
)

//func init() {
//}

type SlackAppMentionEvent struct {
	Type string `json:"type"`
	User string `json:"user"`
	Text string `json:"text"`
	Ts string `json:"ts"`
	Channel string `json:"channel"`
	EventTs string `json:"event_ts"`
}

// TODO Verify token, etc.
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
	fmt.Fprint(w, html.EscapeString(d.Text))
}
