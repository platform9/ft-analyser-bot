package bugsangapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type AllErrors []struct {
	//ID        string `json:"id"`
	Message string `json:"message"`
	//Context   string `json:"context"`
	//Severity  string `json:"severity"`
	EventsURL string `json:"events_url"`
}

type AllEventsOfError []struct {
	//ID           string    `json:"id"`
	URL string `json:"url"`
	/*ProjectURL   string    `json:"project_url"`
	IsFullReport bool      `json:"is_full_report"`
	ErrorID      string    `json:"error_id"`
	ReceivedAt   time.Time `json:"received_at"`
	Exceptions   []struct {
		ErrorClass string `json:"errorClass"`
		Message    string `json:"message"`
	} `json:"exceptions"`
	Severity  string `json:"severity"`
	Context   string `json:"context"`
	Unhandled bool   `json:"unhandled"`
	App       struct {
		ReleaseStage string      `json:"releaseStage"`
		Type         interface{} `json:"type"`
	} `json:"app"`*/
}

type DetailsOFEvent struct {
	//ID   string `json:"id"`
	User struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"user"`
	/*Breadcrumbs []struct {
		Timestamp time.Time `json:"timestamp"`
		Name      string    `json:"name"`
		Type      string    `json:"type"`
		MetaData  struct {
			TargetText     string `json:"targetText"`
			TargetSelector string `json:"targetSelector"`
		} `json:"metaData"`
	} `json:"breadcrumbs"`*/
	Context string `json:"context"`
	//Severity string `json:"severity"`
}

var (
//
//allEvetns = AllEventsOfError{}
//details   = DetailsOFEvent{}
)

func GetAllErrors(userID string) {
	fmt.Println("checking for bugsnag errors for ", userID)
	client := http.Client{}
	//endpoint := fmt.Sprintf("https://api.bugsnag.com/projects/610afe83dd7ab4001516e0d7/errors/?filters[event.since][][value]=7d")
	endpoint := fmt.Sprintf("https://api.bugsnag.com/projects/610afe83dd7ab4001516e0d7/errors/?auth_token=e50ea232-5003-41b9-8315-d76ce9dfa881")
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		fmt.Errorf("Error occured while creating req", err.Error())
	}
	//req.Header.Set("Authorization: token", "e50ea232-5003-41b9-8315-d76ce9dfa881")
	//req.Header.Set("Content-Type", "application/json")
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
	allErrors := AllErrors{}
	err = json.Unmarshal(body, &allErrors)
	if err != nil {
		fmt.Errorf("Unable to unmarshal node info: %w", err)
	}

	for _, eachError := range allErrors {
		url := eachError.EventsURL + "/?auth_token=e50ea232-5003-41b9-8315-d76ce9dfa881"
		allEvetns := GetAllEvents(url)
		if err != nil {
			fmt.Println("failed to get events of error:", eachError.Message)
		}
		for _, eachEvent := range allEvetns {
			url := eachEvent.URL + "/?auth_token=e50ea232-5003-41b9-8315-d76ce9dfa881"
			details := GetDetailsOfEvent(url)
			if err != nil {
				fmt.Println("failed to get event details of err:", err)
			}
			//fmt.Println("Events from : ", url)
			fmt.Println("-----------------------------------------------------------------------------")
			//found := strings.Compare(details.User.ID, "452f36c0612a4e9894c09984cb142d13")
			//if found == 0 {
			fmt.Println("Email:", details.User.Email)
			fmt.Println("ID:", details.User.ID)
			fmt.Println("Erros:", details.Context)
			//}

			//fmt.Println("-----------------------------------------------------------------------------")
		}
	}

}

func GetAllEvents(url string) AllEventsOfError {
	client := http.Client{}
	//endpoint := fmt.Sprintf("https://api.bugsnag.com/projects/610afe83dd7ab4001516e0d7/errors/%s/?filters[event.since][][value]=7d", errorID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Errorf("Error occured while creating req", err.Error())
	}
	//req.Header.Set("Authorization: token", "e50ea232-5003-41b9-8315-d76ce9dfa881")
	//req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Errorf("Error :", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Errorf("unable to send request to bugsnag")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("Unable to read resp body: %w", err)
	}
	allEvents := AllEventsOfError{}
	err = json.Unmarshal(body, &allEvents)
	if err != nil {
		fmt.Errorf("Unable to unmarshal allEvents info: %w", err)
	}
	return allEvents
}

func GetDetailsOfEvent(url string) DetailsOFEvent {
	client := http.Client{}
	//endpoint := fmt.Sprintf("https://api.bugsnag.com/projects/610afe83dd7ab4001516e0d7/errors/%s/?filters[event.since][][value]=7d", errorID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Errorf("Error occured while creating req", err.Error())
	}
	//req.Header.Set("Authorization: token", "e50ea232-5003-41b9-8315-d76ce9dfa881")
	//req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Errorf("Error :", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		//fmt.Println(url)
		//fmt.Errorf("unable to send request to bugsnag")
		fmt.Errorf("unable to send request to bugsnag")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("Unable to read resp body: %w", err)
	}
	details := DetailsOFEvent{}
	err = json.Unmarshal(body, &details)
	if err != nil {
		fmt.Errorf("Unable to unmarshal details of event: %w", err)
	}
	return details
}
