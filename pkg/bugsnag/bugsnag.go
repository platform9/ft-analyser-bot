package bugsnag

// TODO: To be cleaned
import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func GetAllErrors(userID string, email string) ([]UIErrors, error) {
	zap.S().Infof("checking for bugsnag errors for %v", userID)

	client := &http.Client{}

	//bugsnagCreds := config.BugsnagCreds()
	endpoint := fmt.Sprintf(`https://api.bugsnag.com/projects/610afe83dd7ab4001516e0d7/errors?filters[event.since][][value]=7d`)
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

	for i, _ := range allErrors {

		url := fmt.Sprintf(`https://api.bugsnag.com/projects/610afe83dd7ab4001516e0d7/errors/%s/pivots/user.id/values?filters[event.since][][value]=7d&filters[event.since][][type]=eq&filters[user.email][][value]=%s&filters[user.email][][type]=eq`, allErrors[i].ID, email)
		api.baseURL = url
		bugsnagInfo, err = api.getBugsnagInfoAPI()
		if err != nil {
			fmt.Println("Error:", err)
		}

		userError := UserFaced{}
		err = json.Unmarshal(bugsnagInfo, &userError)
		if err != nil {
			zap.S().Errorf("unable to unmarshal error info of user: %w", err)
			//return nil, fmt.Errorf("unable to unmarshal error info of user: %w", err)
		}

		if len(userError) >= 1 {
			uiErrors = append(uiErrors, UIErrors{
				ErrorFaced: allErrors[i].Message,
			})
		}
	}

	return uiErrors, nil
}
