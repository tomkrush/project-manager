package jira

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"tomkrush.com/project-manager/cache"
)

type API struct {
	user, pass, host string
	fileCache        *cache.FileCache
}

func NewAPI(user, pass, host string, fileCache *cache.FileCache) *API {
	return &API{
		user:      user,
		pass:      pass,
		host:      host,
		fileCache: fileCache,
	}
}

func (j *API) Fields() ([]Field, error) {
	queryString := j.host + "/rest/api/3/field"

	body := j.fileCache.Remember(queryString, func() []byte {
		body, err := j.request(queryString)

		if err != nil {
			fmt.Println(err)
		}

		return body
	})

	var data []Field
	err := json.Unmarshal(body, &data)
	if err != nil {
		return []Field{}, err
	}

	return data, nil
}

func (j *API) request(queryString string) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", queryString, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", j.authorizationToken())

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, err
}

func (j *API) MyUser() (User, error) {
	queryString := j.host + "/rest/api/3/myself"

	body := j.fileCache.Remember(queryString, func() []byte {
		body, err := j.request(queryString)

		if err != nil {
			fmt.Println(err)
		}

		return body
	})

	var data User
	err := json.Unmarshal(body, &data)
	if err != nil {
		return User{}, err
	}

	return data, nil
}

func (j *API) authorizationToken() string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(j.user+":"+j.pass))
}

func (j *API) Search(jql string, startAt, maxResults int, fields, expands []string) (SearchResponse, error) {
	params := make(map[string]string)
	cacheKey := ""

	if jql != "" {
		params["jql"] = jql
		cacheKey += jql
	}
	if startAt > 0 {
		params["startAt"] = fmt.Sprintf("%d", startAt)
		cacheKey += params["startAt"]
	}
	if maxResults > 0 {
		params["maxResults"] = fmt.Sprintf("%d", maxResults)
		cacheKey += params["maxResults"]
	}
	if len(fields) > 0 {
		params["fields"] = strings.Join(fields, ",")
		cacheKey += params["fields"]

	}
	if len(expands) > 0 {
		params["expand"] = strings.Join(expands, ",")
		cacheKey += params["expand"]
	}

	queryString := j.host + "/rest/api/3/search?" + mapToQueryString(params)

	body := j.fileCache.Remember(cacheKey, func() []byte {
		body, err := j.request(queryString)

		if err != nil {
			fmt.Println(err)
		}

		return body
	})

	var data SearchResponse
	err := json.Unmarshal(body, &data)
	if err != nil {
		return SearchResponse{}, err
	}

	return data, nil
}

func mapToQueryString(m map[string]string) string {
	query := ""
	for key, value := range m {
		if query != "" {
			query += "&"
		}
		query += fmt.Sprintf("%s=%s", key, value)
	}
	return query
}
