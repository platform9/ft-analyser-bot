package amplitude

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

//Amplitude charts id
var (
	charts = map[string]string{
		"Total User signups":       "erm22ws",
		"User who verified emails": "4hmlqet",
		"Prep-node attempts all":   "8oxoeq7",
		"Prep-node attempts new":   "amqonor",
		"Prep-node errors all":     "vdx7sf1",
		"Prep-node errors new":     "b9vc5pe",
		"Cluster creation all":     "f68bru1",
		"Cluster creation new":     "vyc4vcr",
	}
)

type chartData struct {
	Data struct {
		SeriesCollapsed [][]struct {
			SetID string `json:"setId"`
			Value int    `json:"value"`
		} `json:"seriesCollapsed"`
	} `json:"data"`
}

type WeeklyAnalysis struct {
	Total_User_Signups        int
	User_Who_Verified_Emails  int
	PrepNode_Details          PrepNode
	Cluster_Creation_Attempts Users
}

type PrepNode struct {
	PrepNode_Attempts Users
	PrepNode_Success  Users
	PrepNode_Errors   Users
}

type Users struct {
	New_Users      int
	Existing_Users int
}

// WeeklyMessage formats the fetched amplitude charts info.
func WeeklyMessage() (WeeklyAnalysis, error) {
	result := make(map[string]int)

	// Fetch the count of users performed from each chart
	for chartName, id := range charts {
		client := &http.Client{}
		amplitudeUrl := fmt.Sprintf("https://amplitude.com/api/3/chart/%s/query", id)

		api := api{client, amplitudeUrl}
		chartInfo, err := api.getInfoAPI()
		if err != nil {
			zap.S().Errorf("unable to get chart info from API, error: %v", err.Error())
			return WeeklyAnalysis{}, err
		}

		chart := chartData{}
		err = json.Unmarshal(chartInfo, &chart)
		if err != nil {
			zap.S().Errorf("unable to unmarsha chart info, error: %v", err.Error())
			return WeeklyAnalysis{}, err
		}
		result[chartName] = chart.Data.SeriesCollapsed[0][0].Value
	}

	return generateResult(result), nil
}

func generateResult(result map[string]int) WeeklyAnalysis {
	var weeklyAnalysis WeeklyAnalysis
	weeklyAnalysis.Total_User_Signups = result["Total User signups"]
	weeklyAnalysis.User_Who_Verified_Emails = result["User who verified emails"]
	weeklyAnalysis.PrepNode_Details.PrepNode_Attempts.New_Users = result["Prep-node attempts new"]
	weeklyAnalysis.PrepNode_Details.PrepNode_Attempts.Existing_Users = result["Prep-node attempts all"] - result["Prep-node attempts new"]
	weeklyAnalysis.PrepNode_Details.PrepNode_Success.New_Users = result["Prep-node attempts new"] - result["Prep-node errors new"]
	weeklyAnalysis.PrepNode_Details.PrepNode_Success.Existing_Users = (result["Prep-node attempts all"] - result["Prep-node attempts new"]) - (result["Prep-node errors all"] - result["Prep-node errors new"])
	weeklyAnalysis.PrepNode_Details.PrepNode_Errors.New_Users = result["Prep-node errors new"]
	weeklyAnalysis.PrepNode_Details.PrepNode_Errors.Existing_Users = result["Prep-node errors all"] - result["Prep-node errors new"]
	weeklyAnalysis.Cluster_Creation_Attempts.New_Users = result["Cluster creation new"]
	weeklyAnalysis.Cluster_Creation_Attempts.Existing_Users = result["Cluster creation all"] - result["Cluster creation new"]
	return weeklyAnalysis
}
