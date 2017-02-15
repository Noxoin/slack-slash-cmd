package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "regexp"
)

var slackToken = "FlMln4NEv5xTGc0czw9NKvQ9"
var regex = regexp.MustCompile("[A-Z]+-\\d+")
var jiraPath = "http://jira.indexexchange.com/browse/"

type SlackResponse struct {
    ResponseType string `json:"response_type"`
    Text         string `json:"text"`
    //Attachment   []string `json:"attachment"`
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
    errorJSON, _ := json.Marshal(errorResponse)
    if token != slackToken {
        w.WriteHeader(http.StatusUnauthorized)
        w.Write(errorJSON)
        return
    }
    if command != "/jira" {
        w.WriteHeader(http.StatusBadRequest)
        w.Write(errorJSON)
        return
    }

    response := &SlackResponse {
        ResponseType: "in_channel",
        Text: message,
    }
    tickets := regex.FindAllString(message, -1)
    for i := 0; i < len(tickets); i++ {
        link := fmt.Sprintf("%s: %s%s", tickets[i], jiraPath, tickets[i])
        fmt.Println(link)
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
