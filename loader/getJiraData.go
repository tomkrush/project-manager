package loader

import (
	"fmt"
	"log"
	"math"

	"tomkrush.com/project-manager/jira"
)

type JiraStorage struct {
	Issues jira.Issues
	MyUser jira.User
}

func GetJiraData(jiraAPI *jira.API, jql string) JiraStorage {
	fmt.Println("Sync Jira Data")

	fmt.Println("- Get Tickets")

	jiraStorage := JiraStorage{}
	limit := 50

	numberOfPages, err := getNumberOfPages(jiraAPI, jql, limit)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Number of Requests: %d\n", numberOfPages)

	// Channels for collecting tickets and errors
	ticketsCh := make(chan jira.Issues, numberOfPages)
	errCh := make(chan error, numberOfPages)
	sem := make(chan struct{}, 10) // 5 tokens in semaphore for concurrency control

	// Create go routines to fetch tickets
	for i := 0; i <= numberOfPages; i++ {
		fmt.Printf("Fetching batch %d of %d\n", i, numberOfPages)
		offset := i * limit
		sem <- struct{}{} // Acquire a token
		go func(i int, offset int) {
			defer func() { <-sem }() // Release the token back into the pool
			fetchTickets(jiraAPI, jql, offset, limit, ticketsCh, errCh)
		}(i, offset)
	}

	// Collect tickets and errors
	for i := 0; i <= numberOfPages; i++ {
		select {
		case tickets := <-ticketsCh:
			jiraStorage.Issues = append(jiraStorage.Issues, tickets...)
		case err := <-errCh:
			log.Fatal("Error while fetching tickets: ", err)
		}
	}

	// Close channels
	close(ticketsCh)
	close(errCh)

	fmt.Println("Total Tickets:", len(jiraStorage.Issues))

	fmt.Println("- Get My User")
	myUser, err := jiraAPI.MyUser()

	jiraStorage.MyUser = myUser

	if err != nil {
		log.Fatal(err)
	}

	return jiraStorage
}

func fetchTickets(jiraAPI *jira.API, jiraJQL string, offset int, limit int, out chan jira.Issues, errCh chan error) {
	tickets, err := jiraAPI.Search(jiraJQL, offset, limit, nil, []string{
		"changelog",
	})

	if err != nil {
		errCh <- err
		return
	}
	out <- tickets.Issues
}

func getNumberOfPages(jiraAPI *jira.API, jiraJQL string, limit int) (int, error) {
	tickets, err := jiraAPI.Search(jiraJQL, 0, 1, nil, nil)

	if err != nil {
		log.Fatal(err)
	}

	numberOfPages := int(math.Ceil(float64(tickets.Total) / float64(limit)))

	return numberOfPages, err
}
