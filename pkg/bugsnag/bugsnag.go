package bugsnag

// TODO: To be cleaned
import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/platform9/ft-analyser-bot/pkg/config"
	"go.uber.org/zap"
)

type AllErrors []struct {
	Message   string `json:"message"`
	EventsURL string `json:"events_url"`
}

type AllEventsOfError []struct {
	URL string `json:"url"`
}

type DetailsOFEvent struct {
	User struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"user"`

	Context string `json:"context"`
}

func GetAllErrors(userID string) error {
	fmt.Println("checking for bugsnag errors for ", userID)

	client := &http.Client{}

	bugsnagCreds := config.BugsnagCreds()

	//endpoint := fmt.Sprintf("https://api.bugsnag.com/projects/610afe83dd7ab4001516e0d7/errors/?filters[event.since][][value]=7d")
	//endpoint := fmt.Sprintf("https://api.bugsnag.com/projects/610afe83dd7ab4001516e0d7/errors/?auth_token=e50ea232-5003-41b9-8315-d76ce9dfa881")
	endpoint := fmt.Sprintf("https://api.bugsnag.com/projects/%s/errors/?auth_token=%s", bugsnagCreds.ProjectID, bugsnagCreds.AuthToken)

	api := api{client, endpoint}

	bugsnagInfo, err := api.getBugsnagInfoAPI()
	if err != nil {
		zap.S().Errorf("unable to get bugsnag info from API, error: %v", err.Error())
		return fmt.Errorf("unable to get bugsnag info from API, error: %v", err.Error())
	}

	allErrors := AllErrors{}
	err = json.Unmarshal(bugsnagInfo, &allErrors)
	if err != nil {
		zap.S().Errorf("unable to unmarshal all errors info from bugsnag: %w", err)
		return fmt.Errorf("unable to unmarshal all errors info from bugsnag: %w", err)
	}

	for _, eachError := range allErrors {

		//url := eachError.EventsURL
		url := fmt.Sprintf("%s/?auth_token=%s", eachError.EventsURL, bugsnagCreds.AuthToken)
		allEvents, err := GetAllEvents(url)
		if err != nil {
			zap.S().Errorf("failed to get events of error:", eachError.Message)
			return err
		}

		for _, eachEvent := range allEvents {
			//url := eachEvent.URL
			url := fmt.Sprintf("%s/?auth_token=%s", eachEvent.URL, bugsnagCreds.AuthToken)
			details, err := GetDetailsOfEvent(url)
			if err != nil {
				zap.S().Errorf("failed to get event details of err:", err)
				return err
			}

			if userID == details.User.ID {
				fmt.Println("Email:", details.User.Email)
				fmt.Println("ID:", details.User.ID)
				fmt.Println("Errors:", details.Context)
			}
		}
	}
	return nil
}

// To get list of all events
func GetAllEvents(url string) (AllEventsOfError, error) {
	client := &http.Client{}

	api := api{client, url}

	allEventsResp, err := api.getBugsnagInfoAPI()
	if err != nil {
		return nil, fmt.Errorf("unable to get all events info from API %v, error: %v", url, err.Error())
	}

	allEvents := AllEventsOfError{}
	err = json.Unmarshal(allEventsResp, &allEvents)
	if err != nil {
		return nil, err
	}
	return allEvents, nil
}

// To get details of a event
func GetDetailsOfEvent(url string) (DetailsOFEvent, error) {
	client := &http.Client{}

	api := api{client, url}

	//endpoint := fmt.Sprintf("https://api.bugsnag.com/projects/610afe83dd7ab4001516e0d7/errors/%s/?filters[event.since][][value]=7d", errorID)
	getDetailsInfo, err := api.getBugsnagInfoAPI()
	if err != nil {
		return DetailsOFEvent{}, fmt.Errorf("unable to get details of events info from API %v, error: %v", url, err.Error())
	}

	details := DetailsOFEvent{}
	err = json.Unmarshal(getDetailsInfo, &details)
	if err != nil {
		return DetailsOFEvent{}, err
	}
	return details, nil
}
