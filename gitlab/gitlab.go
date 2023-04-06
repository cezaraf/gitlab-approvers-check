package gitlab

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var httpClient = http.Client{Transport: &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}}

type MergeRequestAuthor struct {
	Id   int8   `json:"id"`
	Name string `json:"name"`
}

type MergeRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type MergeRequestChanges struct {
	MergeRequest
	Author  MergeRequestAuthor   `json:"author"`
	Changes []MergeRequestChange `json:"changes"`
}

type MergeRequestAppprovals struct {
	MergeRequest
	ApprovedBy []MergeRequestApproval `json:"approved_by"`
}

type MergeRequestApproval struct {
	User GitlabUser `json:"user"`
}

type GitlabUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

type MergeRequestChange struct {
	NewPath string `json:"new_path"`
}

type GitlabAccessConfig struct {
	Host           string
	Token          string
	ProjectId      int
	MergeRequestId int
}

func GetMergeRequestChanges(config GitlabAccessConfig) (mergeRequest *MergeRequestChanges, _ error) {

	mergeRequestChangeUrl := fmt.Sprintf("%s/api/v4/projects/%d/merge_requests/%d/changes",
		config.Host, config.ProjectId, config.MergeRequestId)

	httpRequest, err := http.NewRequest("GET", mergeRequestChangeUrl, nil)
	if err != nil {
		return nil, err
	}

	httpRequest.Header = http.Header{
		"PRIVATE-TOKEN": {config.Token},
	}

	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	httpResponseData, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		log.Fatal(err)
	}

	if httpResponse.StatusCode != 200 {
		return nil, fmt.Errorf("the request to url %s returns %s", httpResponse.Request.URL, httpResponse.Status)
	}

	var mergeRequestStruct MergeRequestChanges
	json.Unmarshal(httpResponseData, &mergeRequestStruct)
	return &mergeRequestStruct, nil
}

func GetMergeRequestApprovals(config GitlabAccessConfig) (mergeRequest *MergeRequestAppprovals, _ error) {

	mergeRequestChangeUrl := fmt.Sprintf("%s/api/v4/projects/%d/merge_requests/%d/approvals",
		config.Host, config.ProjectId, config.MergeRequestId)

	httpRequest, err := http.NewRequest("GET", mergeRequestChangeUrl, nil)
	if err != nil {
		return nil, err
	}

	httpRequest.Header = http.Header{
		"PRIVATE-TOKEN": {config.Token},
	}

	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode != 200 {
		return nil, fmt.Errorf("the request to url %s returns %s", httpResponse.Request.URL, httpResponse.Status)
	}

	httpResponseData, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		log.Fatal(err)
	}

	var mergeRequestStruct MergeRequestAppprovals
	json.Unmarshal(httpResponseData, &mergeRequestStruct)
	return &mergeRequestStruct, nil
}

func GetUser(config GitlabAccessConfig, userId int) (user *GitlabUser, _ error) {

	userDataUrl := fmt.Sprintf("%s/api/v4/users/%d", config.Host, userId)

	httpRequest, err := http.NewRequest("GET", userDataUrl, nil)
	if err != nil {
		return nil, err
	}

	httpRequest.Header = http.Header{
		"PRIVATE-TOKEN": {config.Token},
	}

	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode != 200 {
		return nil, fmt.Errorf("the request to url %s returns %s", httpResponse.Request.URL, httpResponse.Status)
	}

	httpResponseData, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		log.Fatal(err)
	}

	var mergeRequestStruct GitlabUser
	json.Unmarshal(httpResponseData, &mergeRequestStruct)
	return &mergeRequestStruct, nil
}
