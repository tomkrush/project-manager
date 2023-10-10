package jira

import (
	"runtime"
	"strings"
	"sync"
)

type Issue struct {
	ID        string    `json:"id"`
	Key       string    `json:"key"`
	Fields    Fields    `json:"fields"`
	Changelog Changelog `json:"changelog"`
}

type Issues []Issue

type SearchResponse struct {
	Issues Issues `json:"issues"`
	Total  int    `json:"total"`
}

func (i Issues) Limit(limit int) Issues {
	if len(i) <= limit {
		return i
	}

	return i[:limit]
}

func (i Issues) Filter(criteria func(issue Issue) bool) Issues {
	var filteredIssues Issues
	var mu sync.Mutex
	var wg sync.WaitGroup

	numberOfCores := runtime.NumCPU()
	chunkSize := len(i) / numberOfCores // Assuming 4 cores; adjust as needed

	filterChunk := func(start, end int) {
		defer wg.Done()
		localFiltered := make(Issues, 0)

		for j := start; j < end; j++ {
			if criteria(i[j]) {
				localFiltered = append(localFiltered, i[j])
			}
		}

		mu.Lock()
		filteredIssues = append(filteredIssues, localFiltered...)
		mu.Unlock()
	}

	for start := 0; start < len(i); start += chunkSize {
		end := start + chunkSize
		if end > len(i) {
			end = len(i)
		}

		wg.Add(1)
		go filterChunk(start, end)
	}

	wg.Wait()

	return filteredIssues
}

func FilterByTextContains(text string) func(issue Issue) bool {
	return func(issue Issue) bool {
		if strings.Contains(issue.Fields.Summary, text) {
			return true
		}

		return strings.Contains(issue.Fields.Description.String(), text)
	}
}

func FilterByStatus(status string) func(issue Issue) bool {
	return func(issue Issue) bool {
		return issue.Fields.Status.Name == status
	}
}

func FilterByAssigneeName(name string) func(issue Issue) bool {
	return func(issue Issue) bool {
		return issue.Fields.Assignee.DisplayName == name
	}
}

func FilterByResolved() func(issue Issue) bool {
	return func(issue Issue) bool {
		return issue.Fields.ResolutionDate != ""
	}
}

func FilterByUnresolved() func(issue Issue) bool {
	return func(issue Issue) bool {
		return issue.Fields.ResolutionDate == ""
	}
}

func (i Issues) GetIssueKeys() []string {
	var keys = []string{}
	for _, issue := range i {
		keys = append(keys, issue.Key)
	}
	return keys
}

type Changelog struct {
	Histories  []History `json:"histories"`
	Total      int       `json:"total"`
	MaxResults int       `json:"maxResults"`
	StartAt    int       `json:"startAt"`
}

type History struct {
	ID      string `json:"id"`
	Author  User   `json:"author"`
	Created string `json:"created"`
	Items   []Item `json:"items"`
}

type Item struct {
	Field      string `json:"field"`
	FieldType  string `json:"fieldtype"`
	From       string `json:"from"`
	FromString string `json:"fromString"`
	To         string `json:"to"`
	ToString   string `json:"toString"`
}

type Fields struct {
	Summary        string       `json:"summary"`
	IssueType      IssueType    `json:"issuetype"`
	Labels         Labels       `json:"labels"`
	Reporter       User         `json:"reporter"`
	Priority       Priority     `json:"Priority"`
	Description    Description  `json:"Description"`
	Created        string       `json:"created"`
	Updated        string       `json:"updated"`
	Assignee       User         `json:"assignee"`
	Status         Status       `json:"status"`
	Creator        User         `json:"creator"`
	Sprint         []Sprint     `json:"customfield_10"`
	Team           Team         `json:"customfield_11"`
	Project        Project      `json:"project"`
	ResolutionDate string       `json:"resolutiondate"`
	FixVersions    []FixVersion `json:"fixVersions"`
}

type Team struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

type Sprint struct {
	Name         string `json:"name"`
	State        string `json:"state"`
	Goal         string `json:"goal"`
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
	CompleteDate string `json:"completeDate"`
}

type FixVersion struct {
	Name         string `json:"name"`
	ReleasedDate string `json:"releasedDate"`
	Released     bool   `json:"released"`
	Archived     bool   `json:"archived"`
}

type Project struct {
	Key        string `json:"key"`
	Name       string `json:"name"`
	AvatarURLs struct {
		URL16x16 string `json:"16x16"`
		URL24x24 string `json:"24x24"`
		URL32x32 string `json:"32x32"`
		URL48x48 string `json:"48x48"`
	} `json:"avatarUrls"`
}

type Status struct {
	Description    string `json:"description"`
	IconUrl        string `json:"iconUrl"`
	Name           string `json:"name"`
	StatusCategory struct {
		Key       string `json:"key"`
		ColorName string `json:"colorName"`
		Name      string `json:"name"`
	} `json:"statusCategory"`
}

type Priority struct {
	IconUrl string `json:"iconUrl"`
	Name    string `json:"name"`
}

type Labels []string

type IssueType struct {
	Name        string `json:"name"`
	IconUrl     string `json:"iconUrl"`
	Description string `json:"description"`
}

type User struct {
	AccountID  string `json:"accountId"`
	Active     bool   `json:"active"`
	AvatarURLs struct {
		URL16x16 string `json:"16x16"`
		URL24x24 string `json:"24x24"`
		URL32x32 string `json:"32x32"`
		URL48x48 string `json:"48x48"`
	} `json:"avatarUrls"`
	DisplayName  string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
}

type Field struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}
