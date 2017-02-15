package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "regexp"
    "strings"
)

var MAX_ATTACHMENTS = 20

var slackToken = "FlMln4NEv5xTGc0czw9NKvQ9"
var regex = regexp.MustCompile("[A-Za-z]+-\\d+")
var jiraPath = "http://jira.indexexchange.com/browse/"

type SlackResponse struct {
    ResponseType string `json:"response_type"`
    Text         string `json:"text"`
    Attachments  []map[string]interface{} `json:"attachments"`
}

func jira(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    token := r.FormValue("token")
    command := r.FormValue("command")
    message := r.FormValue("text")

    // Validation
    errorResponse := &SlackResponse {
        ResponseType: "ephemeral",
        Text: "Sorry, that didn't work. Please try again",
    }
    if token != slackToken {
        errorJSON, _ := json.Marshal(errorResponse)
        w.WriteHeader(http.StatusUnauthorized)
        w.Write(errorJSON)
        return
    }
    if command != "/jira" {
        errorJSON, _ := json.Marshal(errorResponse)
        w.WriteHeader(http.StatusBadRequest)
        w.Write(errorJSON)
        return
    }

    // Processing request
    response := &SlackResponse {
        ResponseType: "in_channel",
    }
    tickets := regex.FindAllString(message, -1)
    for i := 0; i < len(tickets) && i < MAX_ATTACHMENTS; i++ {
        ticket := strings.ToUpper(tickets[i])
        link := fmt.Sprintf("%s: %s%s", ticket, jiraPath, ticket)
        attachment := make(map[string]interface{})
        attachment["text"] =  link
        response.Attachments = append(response.Attachments, attachment)
    }
    responseJSON, _ := json.Marshal(response)
    w.Write(responseJSON)
}

func main() {
    fmt.Printf("Hello, World.\n")
    http.HandleFunc("/jira", jira)
    http.ListenAndServe(":8080", nil)
}
