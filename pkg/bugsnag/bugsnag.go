package bugsnag

// TODO: To be cleaned
import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/platform9/ft-analyser-bot/pkg/config"
	"go.uber.org/zap"
)

type AllErrors []struct {
	URL string `json:"url"`
}

type ErrorDetails struct {
	EventsURL string `json:"events_url"`
}

type AllEventsOfError []struct {
	URL string `json:"url"`
}

type UIErrors struct {
	UserID     string
	ErrorFaced string
	Time       time.Time
}

type EventDetails struct {
	ID           string    `json:"id"`
	URL          string    `json:"url"`
	ProjectURL   string    `json:"project_url"`
	IsFullReport bool      `json:"is_full_report"`
	ErrorID      string    `json:"error_id"`
	ReceivedAt   time.Time `json:"received_at"`
	Exceptions   []struct {
		ErrorClass string `json:"error_class"`
		Message    string `json:"message"`
		Type       string `json:"type"`
		Stacktrace []struct {
			ColumnNumber      int         `json:"column_number"`
			InProject         interface{} `json:"in_project"`
			LineNumber        int         `json:"line_number"`
			Method            string      `json:"method"`
			File              string      `json:"file"`
			Type              interface{} `json:"type"`
			Code              interface{} `json:"code"`
			CodeFile          interface{} `json:"code_file"`
			AddressOffset     interface{} `json:"address_offset"`
			MachoUUID         interface{} `json:"macho_uuid"`
			RelativeAddress   interface{} `json:"relative_address"`
			FrameAddress      interface{} `json:"frame_address"`
			SourceControlLink interface{} `json:"source_control_link"`
			SourceControlName string      `json:"source_control_name"`
		} `json:"stacktrace"`
		Registers interface{} `json:"registers"`
	} `json:"exceptions"`
	Threads  interface{} `json:"threads"`
	MetaData struct {
		App struct {
			CustomerTier string `json:"customerTier"`
			DuVersion    string `json:"duVersion"`
		} `json:"App"`
		Device struct {
			DeviceMemory        int    `json:"deviceMemory"`
			HardwareConcurrency int    `json:"hardwareConcurrency"`
			AppCodeName         string `json:"appCodeName"`
			AppName             string `json:"appName"`
			AppVersion          string `json:"appVersion"`
			CookieEnabled       bool   `json:"cookieEnabled"`
			MaxTouchPoints      int    `json:"maxTouchPoints"`
			OnLine              bool   `json:"onLine"`
			Platform            string `json:"platform"`
			Product             string `json:"product"`
			ProductSub          string `json:"productSub"`
			Vendor              string `json:"vendor"`
			VendorSub           string `json:"vendorSub"`
			Webdriver           bool   `json:"webdriver"`
		} `json:"Device"`
		ResponseConfig struct {
			URL     string `json:"url"`
			Method  string `json:"method"`
			Headers struct {
				Accept        string `json:"Accept"`
				Authorization string `json:"Authorization"`
				XAuthToken    string `json:"X-Auth-Token"`
			} `json:"headers"`
			TransformRequest  []interface{} `json:"transformRequest"`
			TransformResponse []interface{} `json:"transformResponse"`
			Timeout           int           `json:"timeout"`
			XSRFCookieName    string        `json:"xsrfCookieName"`
			XSRFHeaderName    string        `json:"xsrfHeaderName"`
			MaxContentLength  int           `json:"maxContentLength"`
			MaxBodyLength     int           `json:"maxBodyLength"`
		} `json:"Response Config"`
		Request struct {
			BSU string `json:"BS~~U"`
			BSM string `json:"BS~~M"`
			BSS bool   `json:"BS~~S"`
		} `json:"Request"`
		ResponseData struct {
		} `json:"Response Data"`
		Stacktrace struct {
		} `json:"Stacktrace"`
	} `json:"metaData"`
	Request struct {
		URL      string `json:"url"`
		ClientIP string `json:"clientIp"`
		Headers  struct {
		} `json:"headers"`
	} `json:"request"`
	App struct {
		Version      string `json:"version"`
		ReleaseStage string `json:"releaseStage"`
		Duration     int    `json:"duration"`
	} `json:"app"`
	Device struct {
		ID             string    `json:"id"`
		OsName         string    `json:"osName"`
		BrowserName    string    `json:"browserName"`
		BrowserVersion string    `json:"browserVersion"`
		Orientation    string    `json:"orientation"`
		Locale         string    `json:"locale"`
		Time           time.Time `json:"time"`
	} `json:"device"`
	User struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"user"`
	Breadcrumbs []struct {
		Timestamp time.Time `json:"timestamp"`
		Name      string    `json:"name"`
		Type      string    `json:"type"`
		MetaData  struct {
			Status  int    `json:"status"`
			Request string `json:"request"`
		} `json:"metaData"`
	} `json:"breadcrumbs"`
	Context      string        `json:"context"`
	Severity     string        `json:"severity"`
	Unhandled    bool          `json:"unhandled"`
	FeatureFlags []interface{} `json:"feature_flags"`
}

func GetAllErrors(userID string) error {
	fmt.Println("checking for bugsnag errors for ", userID)

	client := &http.Client{}

	bugsnagCreds := config.BugsnagCreds()

	//endpoint := fmt.Sprintf("https://api.bugsnag.com/projects/PROJECTID/errors/?filters[event.since][][value]=7d")
	//endpoint := fmt.Sprintf("https://api.bugsnag.com/projects/PROJECTID/errors/?auth_token=AUTHTOKEN")
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
	var uiErrors []UIErrors

	for _, eachError := range allErrors {

		//url := eachError.EventsURL
		getErrorDetailsUrl := fmt.Sprintf("%s/?auth_token=%s", eachError.URL, bugsnagCreds.AuthToken)
		errorDetails, err := GetErrorDetails(getErrorDetailsUrl)
		if err != nil {
			zap.S().Errorf("failed to error details:", err)
			return err
		}

		//url := eachEvent.URL
		getAllEventsUrl := fmt.Sprintf("%s/?auth_token=%s", errorDetails.EventsURL, bugsnagCreds.AuthToken)
		allEventsOfError, err := GetAllEventsOfError(getAllEventsUrl)
		if err != nil {
			zap.S().Errorf("failed to get all events of error:", err)
			return err
		}
		var newEventsUrl string
		if allEventsOfError[0].URL != "" {
			newEventsUrl = fmt.Sprintf("%s/?auth_token=%s", allEventsOfError[0].URL, bugsnagCreds.AuthToken)

		} else {
			fmt.Println("Event URL not found")
		}
		eventDetails, err := GetEventDetails(newEventsUrl)
		if err != nil {
			zap.S().Errorf("failed to get event details of err:", err)
			return err
		}

		//TODO: This logic to be verified.
		if userID == eventDetails.User.ID {
			uiErrors = append(uiErrors, UIErrors{UserID: eventDetails.User.ID,
				ErrorFaced: eventDetails.Context,
				Time:       eventDetails.ReceivedAt})
			/*(fmt.Println("UserID:", eventDetails.User.ID)
			fmt.Println("Error faced:", eventDetails.Context)
			fmt.Println("Time:", eventDetails.ReceivedAt)*/
		}
	}

	return nil
}

//

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
