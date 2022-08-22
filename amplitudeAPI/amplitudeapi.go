package amplitudeapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

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

type User struct {
	Matches []struct {
		AmplitudeID int64  `json:"amplitude_id"`
		Country     string `json:"country"`
		LastSeen    string `json:"last_seen"`
	} `json:"matches"`
}

type ChartData struct {
	Data struct {
		SeriesCollapsed [][]struct {
			SetID string `json:"setId"`
			Value int    `json:"value"`
		} `json:"seriesCollapsed"`
	} `json:"data"`
}

func WeeklyMessage() {

	result := make(map[string]int)

	for chartName, id := range charts {
		client := http.Client{}
		endpoint := fmt.Sprintf("https://amplitude.com/api/3/chart/%s/query", id)
		req, err := http.NewRequest(http.MethodGet, endpoint, nil)
		if err != nil {
			fmt.Errorf("Error occured while creating req", err.Error())
		}
		req.Header.Add("Accept", "application/json")
		req.SetBasicAuth("userID", "Password")
		resp, err := client.Do(req)
		if err != nil {
			fmt.Errorf("Error :", err.Error())
		}
		if resp.StatusCode != http.StatusOK {
			fmt.Println("unable to send request to amplitude")
		}
		chart := ChartData{}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Errorf("Unable to read resp body of node info: %w", err)
		}
		err = json.Unmarshal(body, &chart)
		if err != nil {
			fmt.Errorf("Unable to unmarshal node info: %w", err)
		}
		if chartName == "Total User signups" {
			result["Total User signups"] = chart.Data.SeriesCollapsed[0][0].Value
		} else if chartName == "User who verified emails" {
			result["User who verified emails"] = chart.Data.SeriesCollapsed[0][0].Value
		} else if chartName == "Prep-node attempts all" {
			result["Prep-node attempts all"] = chart.Data.SeriesCollapsed[0][0].Value
		} else if chartName == "Prep-node attempts new" {
			result["Prep-node attempts new"] = chart.Data.SeriesCollapsed[0][0].Value
		} else if chartName == "Prep-node errors all" {
			result["Prep-node errors all"] = chart.Data.SeriesCollapsed[0][0].Value
		} else if chartName == "Prep-node errors new" {
			result["Prep-node errors new"] = chart.Data.SeriesCollapsed[0][0].Value
		} else if chartName == "Cluster creation all" {
			result["Cluster creation all"] = chart.Data.SeriesCollapsed[0][0].Value
		} else if chartName == "Cluster creation new" {
			result["Cluster creation new"] = chart.Data.SeriesCollapsed[0][0].Value
		}
	}

	fmt.Println("Total User signups :", result["Total User signups"])
	fmt.Println("User who verified emails :", result["User who verified emails"])
	fmt.Println("Prep-node details")
	fmt.Println("	Prep-node attempts")
	fmt.Println("		New users:", result["Prep-node attempts new"])
	fmt.Println("		Existing users:", result["Prep-node attempts all"]-result["Prep-node attempts new"])
	fmt.Println("	Prep-node successes")
	fmt.Println("		New users:", result["Prep-node attempts new"]-result["Prep-node errors new"])
	fmt.Println("		Existing users:", (result["Prep-node attempts all"]-result["Prep-node attempts new"])-(result["Prep-node errors all"]-result["Prep-node errors new"]))
	fmt.Println("	Prep-node errors")
	fmt.Println("		New users:", result["Prep-node errors new"])
	fmt.Println("		Existing users:", result["Prep-node errors all"]-result["Prep-node errors new"])
	fmt.Println("Cluster creation attempts")
	fmt.Println("	New users:", result["Cluster creation new"])
	fmt.Println("	Existing users:", result["Cluster creation all"]-result["Cluster creation new"])

}

type DUdetails struct {
	Details struct {
		Fqdn     string `json:"fqdn"`
		Metadata struct {
			ActiveHosts string `json:"active_hosts"`
			Bits        string `json:"bits"`
			TotalHosts  string `json:"total_hosts"`
			Version     string `json:"version"`
		} `json:"metadata"`
		TaskState string `json:"task_state"`
	} `json:"details"`
}

type AmplitudeData struct {
	Events []struct {
		EventType          string      `json:"event_type"`
		AmplitudeEventType interface{} `json:"amplitude_event_type"`
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

func NPS_Score_Analysis(userID string) {

	id, err := getAmplitudeID(userID)
	if err != nil {
		fmt.Println("failed to get amplitude id")
	}
	_, err = printUserData(id)
	if err != nil {
		fmt.Println("error getting user data :", err)
	}

}

func getAmplitudeID(userID string) (int64, error) {
	client := http.Client{}
	endpoint := fmt.Sprintf("https://amplitude.com/api/2/usersearch?user=%s", userID)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return 0, fmt.Errorf("Error occured while creating req", err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth("userID", "Password")
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("Error :", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unable to send request to amplitude")
	}
	user := User{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("Unable to read resp body of node info: %w", err)
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		return 0, fmt.Errorf("Unable to unmarshal node info: %w", err)
	}
	return user.Matches[0].AmplitudeID, nil
}

func printUserData(AmplitudeUserID int64) (string, error) {
	client := http.Client{}
	endpoint := fmt.Sprintf("https://amplitude.com/api/2/useractivity?user=%d", AmplitudeUserID)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("Error occured while creating req", err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth("userID", "Password")
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error :", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unable to send request to amplitude")
	}
	Data := AmplitudeData{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Unable to read resp body of node info: %w", err)
	}
	err = json.Unmarshal(body, &Data)
	if err != nil {
		return "", fmt.Errorf("Unable to unmarshal node info: %w", err)
	}

	fmt.Println("User is from ", Data.UserData.Country)
	fmt.Println("First seen on ", Data.UserData.FirstUsed)
	fmt.Println("Last seen on ", Data.UserData.LastUsed)
	err = hostDetails(Data.UserData.Properties.DuFqdn)
	if err != nil {
		fmt.Println("error getting host details")
	}
	fmt.Printf("Some of the user activites are (")
	for i := 0; i <= 9; i++ {
		fmt.Printf("%s", Data.Events[i].EventType+", ")
	}
	fmt.Printf("etc)\n")
	return Data.UserData.Properties.DuFqdn, nil
}

func hostDetails(fqdn string) error {
	client := http.Client{}
	endpoint := fmt.Sprintf("https://bork-prod.platform9.horse/api/v1/regions/%s", fqdn)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("Error occured while creating req", err.Error())
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error :", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to send request to amplitude")
	}
	Data := DUdetails{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("Unable to read resp body of node info: %w", err)
	}
	err = json.Unmarshal(body, &Data)
	if err != nil {
		fmt.Errorf("Unable to unmarshal node info: %w", err)
	}
	fmt.Println("Total hosts attached to DU are ", Data.Details.Metadata.TotalHosts)
	fmt.Println("Active hosts are ", Data.Details.Metadata.ActiveHosts)
	return nil
}
