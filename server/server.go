package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"tomkrush.com/project-manager/jira"
	"tomkrush.com/project-manager/loader"
)

//go:embed static/*
var content embed.FS

func Start(data loader.JiraStorage) {
	indexRoute()
	ticketsRoute(data)
	myselfRoute(data)

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}

func myselfRoute(data loader.JiraStorage) {
	http.HandleFunc("/myself", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data.MyUser)
	})
}

type IssuesResponse struct {
	Issues jira.Issues `json:"issues"`
	Total  int         `json:"total"`
}

type IssuesKeysResponse struct {
	IssueKeys []string `json:"issueKeys"`
	Total     int      `json:"total"`
}

func ticketsRoute(data loader.JiraStorage) {
	http.HandleFunc("/issues", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		issues, total := filterIssuesUsingRequest(data.Issues, r)

		response := IssuesResponse{
			Issues: issues,
			Total:  total,
		}

		json.NewEncoder(w).Encode(response)
	})

	http.HandleFunc("/issueKeys", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		issues, total := filterIssuesUsingRequest(data.Issues, r)

		response := IssuesKeysResponse{
			IssueKeys: issues.GetIssueKeys(),
			Total:     total,
		}

		json.NewEncoder(w).Encode(response)
	})
}

func filterIssuesUsingRequest(issues jira.Issues, r *http.Request) (jira.Issues, int) {
	status := r.URL.Query().Get("status")

	if status != "" {
		issues = issues.Filter(jira.FilterByStatus(status))
	}

	text := r.URL.Query().Get("text")

	if text != "" {
		issues = issues.Filter(jira.FilterByTextContains(text))
	}

	total := len(issues)

	// Get the "limit" query parameter from the URL and set a default value
	limitParam := r.URL.Query().Get("limit")
	defaultLimit := 10

	// Check if the "limit" parameter is provided and parse it as an integer
	if limitParam != "" {
		parsedLimit, err := strconv.Atoi(limitParam)
		if err == nil {
			defaultLimit = parsedLimit
		}
	}

	issues = issues.Limit(defaultLimit)
	return issues, total
}

func indexRoute() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		data, err := content.ReadFile("static/index.html")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Write(data)
		}
	})

	http.HandleFunc("/assets/app.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")

		data, err := content.ReadFile("static/app.js")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Write(data)
		}
	})
}
