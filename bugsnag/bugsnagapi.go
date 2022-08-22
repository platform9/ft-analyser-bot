package bugsangapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type AllErrors []struct {
	ID       string `json:"id"`
	Message  string `json:"message"`
	Context  string `json:"context"`
	Severity string `json:"severity"`
}

var (
	allErrors = AllErrors{}
)

func GetAllErrors() {
	client := http.Client{}
	endpoint := fmt.Sprintf("https://api.bugsnag.com/projects/610afe83dd7ab4001516e0d7/errors")
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		fmt.Errorf("Error occured while creating req", err.Error())
	}
	req.Header.Set("Authorization: token", "")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Errorf("Error :", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("unable to send request to bugsnag")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("Unable to read resp body of node info: %w", err)
	}
	err = json.Unmarshal(body, &allErrors)
	if err != nil {
		fmt.Errorf("Unable to unmarshal node info: %w", err)
	}
	fmt.Println("Today's errors faced by all users")
	for _, err := range allErrors {
		if err.Severity == "error" {
			fmt.Println(err.Message)
		}
	}
}
