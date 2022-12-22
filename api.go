package jira

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"encoding/json"
)

type TransitionBody struct {
	Id int `json:"id"`
}

type MoveTicketBody struct {
	Transition TransitionBody `json:"transition"`
}

func MoveTicketStatus(issueKey string, status int) error {
	client := &http.Client{}
	url := fmt.Sprintf("%s/rest/api/2/issue/%s/transitions", os.Getenv("JIRA_URL"), issueKey)
	body := MoveTicketBody{Transition: TransitionBody{Id: status}}
	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(os.Getenv("JIRA_EMAIL"), os.Getenv("JIRA_API_TOKEN"))
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("Status updated successfully")
	}

	return nil
}

type AssignTicketBody struct {
	AccountID string `json:"accountId"`
}

type GetTicketResponse struct {
	Fields struct {
		Title       string `json:"summary"`
		Description string `json:"description"`
	} `json:"fields"`
}

func GetTicket(issueKey string) (*GetTicketResponse, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/rest/api/2/issue/%s", os.Getenv("JIRA_URL"), issueKey)
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return &GetTicketResponse{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(os.Getenv("JIRA_EMAIL"), os.Getenv("JIRA_API_TOKEN"))
	resp, err := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)

	if err != nil {
		return &GetTicketResponse{}, err
	}

	var respJson = &GetTicketResponse{}
	json.Unmarshal(body, &respJson)

	return respJson, nil
}

func AutoAssignTicket(issueKey string) error {
	client := &http.Client{}
	url := fmt.Sprintf("%s/rest/api/2/issue/%s/assignee", os.Getenv("JIRA_URL"), issueKey)
	body := AssignTicketBody{AccountID: os.Getenv("JIRA_PROFILE_ID")}
	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(os.Getenv("JIRA_EMAIL"), os.Getenv("JIRA_API_TOKEN"))
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode == 204 {
		fmt.Println("Assignee updated successfully")
	}

	return nil
}
