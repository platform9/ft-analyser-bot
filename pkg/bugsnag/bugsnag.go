package bugsnag

// TODO: To be cleaned
import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/platform9/ft-analyser-bot/pkg/config"
	"go.uber.org/zap"
)

func GetAllErrors(userID string) ([]UIErrors, error) {
	zap.S().Infof("checking for bugsnag errors for %v", userID)

	client := &http.Client{}

	bugsnagCreds := config.BugsnagCreds()

	endpoint := fmt.Sprintf("https://api.bugsnag.com/projects/%s/errors/?auth_token=%s", bugsnagCreds.ProjectID, bugsnagCreds.AuthToken)

	api := api{client, endpoint}

	bugsnagInfo, err := api.getBugsnagInfoAPI()
	if err != nil {
		zap.S().Errorf("unable to get bugsnag info from API, error: %v", err.Error())
		return nil, fmt.Errorf("unable to get bugsnag info from API, error: %v", err.Error())
	}

	allErrors := AllErrors{}
	err = json.Unmarshal(bugsnagInfo, &allErrors)
	if err != nil {
		zap.S().Errorf("unable to unmarshal all errors info from bugsnag: %w", err)
		return nil, fmt.Errorf("unable to unmarshal all errors info from bugsnag: %w", err)
	}
	var uiErrors []UIErrors

	for _, eachError := range allErrors {

		//url := eachError.EventsURL
		getErrorDetailsUrl := fmt.Sprintf("%s/?auth_token=%s", eachError.URL, bugsnagCreds.AuthToken)
		errorDetails, err := GetErrorDetails(getErrorDetailsUrl)
		if err != nil {
			zap.S().Errorf("failed to error details:", err)
			return nil, err
		}

		//url := eachEvent.URL
		if errorDetails.EventsURL == "" {
			zap.S().Debug("Events URL not found")
			continue
		}

		getAllEventsUrl := fmt.Sprintf("%s/?auth_token=%s", errorDetails.EventsURL, bugsnagCreds.AuthToken)
		allEventsOfError, err := GetAllEventsOfError(getAllEventsUrl)
		if err != nil {
			zap.S().Errorf("failed to get all events of error:", err)
			return nil, err
		}

		if allEventsOfError[0].URL == "" {
			zap.S().Debug("Event URL not found")
			continue
		}

		newEventsUrl := fmt.Sprintf("%s/?auth_token=%s", allEventsOfError[0].URL, bugsnagCreds.AuthToken)
		eventDetails, err := GetEventDetails(newEventsUrl)
		if err != nil {
			zap.S().Errorf("failed to get event details of err:", err)
			return nil, err
		}

		if userID == eventDetails.User.ID {
			uiErrors = append(uiErrors, UIErrors{UserID: eventDetails.User.ID,
				ErrorFaced: eventDetails.Context,
				Time:       eventDetails.ReceivedAt})
		}
	}

	return uiErrors, nil
}

func GetErrorDetails(url string) (ErrorDetails, error) {
	client := &http.Client{}

	api := api{client, url}

	getErrorDetailsInfo, err := api.getBugsnagInfoAPI()
	if err != nil {
		return ErrorDetails{}, fmt.Errorf("unable to get error details info from API %v, error: %v", url, err.Error())
	}

	errorDetails := ErrorDetails{}
	err = json.Unmarshal(getErrorDetailsInfo, &errorDetails)
	if err != nil {
		fmt.Errorf("unable to unmarshal errorDetails: %w", err)
		return ErrorDetails{}, err
	}
	return errorDetails, nil
}

func GetAllEventsOfError(url string) (AllEventsOfError, error) {
	client := &http.Client{}

	api := api{client, url}

	getAllEventsOfErrorsInfo, err := api.getBugsnagInfoAPI()
	if err != nil {
		return AllEventsOfError{}, fmt.Errorf("unable to get error details info from API %v, error: %v", url, err.Error())
	}

	allEventsOfError := AllEventsOfError{}
	err = json.Unmarshal(getAllEventsOfErrorsInfo, &allEventsOfError)
	if err != nil {
		fmt.Errorf("unable to unmarshal allEventsOfError: %w", err)
		return AllEventsOfError{}, err
	}
	return allEventsOfError, nil
}

func GetEventDetails(url string) (EventDetails, error) {
	client := &http.Client{}
	api := api{client, url}

	getEventDetailsInfo, err := api.getBugsnagInfoAPI()
	if err != nil {
		return EventDetails{}, fmt.Errorf("unable to get error details info from API %v, error: %v", url, err.Error())
	}

	eventDetails := EventDetails{}
	err = json.Unmarshal(getEventDetailsInfo, &eventDetails)
	if err != nil {
		fmt.Errorf("Unable to unmarshal eventDetails: %w", err)
		return EventDetails{}, err
	}
	return eventDetails, nil
}
