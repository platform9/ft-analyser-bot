package amplitude

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/platform9/ft-analyser-bot/pkg/bugsnag"
	"go.uber.org/zap"
)

type user struct {
	Matches []struct {
		AmplitudeID int64  `json:"amplitude_id"`
		Country     string `json:"country"`
		LastSeen    string `json:"last_seen"`
	} `json:"matches"`
}

type duDetails struct {
	Details struct {
		Fqdn     string `json:"fqdn"`
		Name     string `json:"customer_shortname"`
		Metadata struct {
			ActiveHosts string `json:"active_hosts"`
			Bits        string `json:"bits"`
			TotalHosts  string `json:"total_hosts"`
			Version     string `json:"version"`
		} `json:"metadata"`
		Cluster   string `json:"cluster"`
		TaskState string `json:"task_state"`
	} `json:"details"`
}

type amplitudeData struct {
	Events []struct {
		EventType string `json:"event_type"`
	} `json:"events"`
	UserData struct {
		FirstSeen  int64 `json:"firstSeen"`
		LastSeen   int64 `json:"lastSeen"`
		Properties struct {
			DuFqdn string `json:"du_fqdn"`
		} `json:"properties"`
		Version   string `json:"version"`
		Country   string `json:"country"`
		FirstUsed string `json:"first_used"`
		LastUsed  string `json:"last_used"`
	} `json:"userData"`
}

type NPSAnalysis struct {
	UserCountry    string
	FirstSeen      string
	LastSeen       string
	UserActivities []string
	CLIEvents      CLI
	HostDetails    DUDetails
	UIErrors       []bugsnag.UIErrors
}

type CLI struct {
	Prepnode struct {
		Success int
		Failure int
	}
	ChecknodeSuccess int
	ClusterCreation  struct {
		Success int
		Delete  int
	}
}

type DUDetails struct {
	FQDN        string
	BorkCluster string
	HostCount   string
	ActiveHosts string
}

/* NpsScoreAnalysis generates a detailed report with
   user info, events and error analysis like UI errors, pf9cli error.
*/
func NpsScoreAnalysis(userID string) (NPSAnalysis, error) {
	// Fetch the amplitude ID
	id, err := getAmplitudeID(userID)
	if err != nil {
		zap.S().Errorf("failed to get amplitude id, error: %v", err.Error())
		return NPSAnalysis{}, err
	}

	//TODO: To be removed, i.e we don't print it just generate the info.
	npsAnalysis, err := printUserData(id)
	if err != nil {
		zap.S().Errorf("failed to print user data, error: %v", err.Error())
		return NPSAnalysis{}, err
	}

	npsAnalysis.UIErrors, err = bugsnag.GetAllErrors(userID)
	if err != nil {
		zap.S().Errorf("failed to fetch UI Errors, error: %v", err.Error())
		//return NPSAnalysis{}, err
	}

	return npsAnalysis, nil
}

// To get amplitude ID
func getAmplitudeID(userID string) (int64, error) {
	client := &http.Client{}

	userSearchUrl := fmt.Sprintf("https://amplitude.com/api/2/usersearch?user=%s", userID)
	api := api{client, userSearchUrl}

	userInfo, err := api.getInfoAPI()
	if err != nil {
		return 0, fmt.Errorf("unable to get user search info from API, error: %v", err.Error())
	}

	user := user{}
	err = json.Unmarshal(userInfo, &user)
	if err != nil {
		return 0, fmt.Errorf("unable to unmarshal user search info, error: %v", err.Error())
	}
	return user.Matches[0].AmplitudeID, nil
}

// To print the user full activity data.
func printUserData(AmplitudeUserID int64) (NPSAnalysis, error) {
	client := &http.Client{}
	userActivityUrl := fmt.Sprintf("https://amplitude.com/api/2/useractivity?user=%d", AmplitudeUserID)
	api := api{client, userActivityUrl}

	userActivityInfo, err := api.getInfoAPI()
	if err != nil {
		return NPSAnalysis{}, fmt.Errorf("unable to get user activity info from API, error:", err.Error())
	}

	userData := amplitudeData{}

	err = json.Unmarshal(userActivityInfo, &userData)
	if err != nil {
		return NPSAnalysis{}, fmt.Errorf("unable to unmarshal user activity info, error: %v", err.Error())
	}

	var npsAnalysis NPSAnalysis
	npsAnalysis.UserCountry = userData.UserData.Country
	npsAnalysis.FirstSeen = userData.UserData.FirstUsed
	npsAnalysis.LastSeen = userData.UserData.LastUsed
	cliDetails := cliDetailsInfo(userData.Events)
	npsAnalysis.CLIEvents = cliDetails

	npsAnalysis.UserActivities = removeDuplicates(userData.Events)[:9]

	//Fetch the DU host details for given fqdn
	hostDetails, err := hostDetails(userData.UserData.Properties.DuFqdn)
	if err != nil {
		return NPSAnalysis{}, fmt.Errorf("error getting host details, error: %v", err.Error())
	}
	npsAnalysis.HostDetails = hostDetails
	return npsAnalysis, nil
}

// To fetch events count
func cliDetailsInfo(strList eventType) CLI {
	var cliDetails CLI
	for i := 0; i < len(strList); i++ {
		switch strList[i].EventType {
		case "qbert: cluster creation success":
			cliDetails.ClusterCreation.Success += 1
		case "qbert: cluster deletion success":
			cliDetails.ClusterCreation.Delete += 1
		case "Prep-node : Successful":
			cliDetails.Prepnode.Success += 1
		case "CheckNode complete":
			cliDetails.ChecknodeSuccess += 1
		case "Prep-node : ERROR":
			cliDetails.Prepnode.Failure += 1
		}
	}
	return cliDetails
}

// To fetch the host details using bork apis.
func hostDetails(fqdn string) (DUDetails, error) {
	client := &http.Client{}
	borkRegionUrl := fmt.Sprintf("https://bork-prod.platform9.horse/api/v1/regions/%s", fqdn)
	//zap.S().Debugf("Bork URL is: %s", borkRegionUrl)
	api := api{client, borkRegionUrl}

	regionInfo, err := api.getInfoAPI()
	if err != nil {
		return DUDetails{}, fmt.Errorf("unable to get DU info from API, error: %v", err.Error())
	}
	duData := duDetails{}
	//zap.S().Debugf("\nduData: %s\n", string(regionInfo))
	err = json.Unmarshal(regionInfo, &duData)
	if err != nil {
		return DUDetails{}, fmt.Errorf("unable to unmarshal bork info, error: %v", err.Error())
	}

	return DUDetails{FQDN: duData.Details.Fqdn,
		BorkCluster: duData.Details.Cluster,
		HostCount:   duData.Details.Metadata.TotalHosts,
		ActiveHosts: duData.Details.Metadata.ActiveHosts}, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

type eventType []struct {
	EventType string "json:\"event_type\""
}

// To remove duplicates from amplitude events list
func removeDuplicates(strList eventType) []string {
	list := []string{}
	for i := 0; i < len(strList); i++ {
		if contains(list, strList[i].EventType) == false {
			list = append(list, strList[i].EventType)
		}
	}
	return list
}
