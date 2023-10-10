package main

import (
	"flag"
	"time"

	"tomkrush.com/project-manager/cache"
	"tomkrush.com/project-manager/jira"
	"tomkrush.com/project-manager/loader"
	"tomkrush.com/project-manager/server"
)

func main() {
	// Declare string pointers for each flag with idiomatic names
	jiraHost := flag.String("host", "", "The JIRA host URL.")
	jiraUser := flag.String("user", "", "The JIRA username.")
	jiraToken := flag.String("token", "", "The JIRA access token.")
	jiraJQL := flag.String("jql", "", "The JIRA JQL query.")

	// Parse the flags
	flag.Parse()

	// Validate the flags
	if *jiraHost == "" || *jiraUser == "" || *jiraToken == "" || *jiraJQL == "" {
		flag.PrintDefaults()
		return
	}

	cache := cache.NewFileCache("/tmp/project-manager", 7*24*time.Hour)
	jiraAPI := jira.NewAPI(*jiraUser, *jiraToken, *jiraHost, cache)

	data := loader.GetJiraData(jiraAPI, *jiraJQL)

	server.Start(data)
}
