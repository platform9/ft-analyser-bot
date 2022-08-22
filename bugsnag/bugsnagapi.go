package bugsangapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
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

func GetAllErrors(userID string) {
	fmt.Println("checking for bugsnag errors for ", userID)
	client := http.Client{}
	//endpoint := fmt.Sprintf("https://api.bugsnag.com/projects/610afe83dd7ab4001516e0d7/errors/?filters[event.since][][value]=7d/?auth_token=e50ea232-5003-41b9-8315-d76ce9dfa881")
	// TO DO : get error of last 7 days (currently getting only present day's errors)
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
		url := eachError.URL + "/?auth_token=e50ea232-5003-41b9-8315-d76ce9dfa881"
		errorDetails := GetErrorDetails(url)
		if err != nil {
			fmt.Println("failed to get error details")
		}
		allEventsOfError := GetAllEventsOfError(errorDetails.EventsURL + "/?auth_token=e50ea232-5003-41b9-8315-d76ce9dfa881")
		var url2 string
		if allEventsOfError[0].URL != "" {
			url2 = allEventsOfError[0].URL + "/?auth_token=e50ea232-5003-41b9-8315-d76ce9dfa881"
		} else {
			fmt.Println("Event URL not found")
		}
		eventDetails := GetEventDetails(url2)

		if userID == eventDetails.User.ID {
			fmt.Println("UserID:", eventDetails.User.ID)
			fmt.Println("Error faced:", eventDetails.Context)
			fmt.Println("Time:", eventDetails.ReceivedAt)
		}

	}

}

func GetErrorDetails(url string) ErrorDetails {
	client := http.Client{}

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
	errorDetails := ErrorDetails{}
	err = json.Unmarshal(body, &errorDetails)
	if err != nil {
		fmt.Errorf("Unable to unmarshal errorDetails: %w", err)
	}
	return errorDetails
}

func GetAllEventsOfError(url string) AllEventsOfError {
	client := http.Client{}

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
	allEventsOfError := AllEventsOfError{}
	err = json.Unmarshal(body, &allEventsOfError)
	if err != nil {
		fmt.Errorf("Unable to unmarshal allEventsOfError: %w", err)
	}
	return allEventsOfError
}

func GetEventDetails(url string) EventDetails {
	client := http.Client{}

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
	eventDetails := EventDetails{}
	err = json.Unmarshal(body, &eventDetails)
	if err != nil {
		fmt.Errorf("Unable to unmarshal eventDetails: %w", err)
	}
	return eventDetails
}
