package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	amplitudeapi "github.com/platform9/ft-analyser-bot/pkg/amplitude"
	"go.uber.org/zap"
)

// New returns new API router for ft-analyser
func New() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/weeklyanalysis", weeklyFTAnalysis).Methods("GET")
	r.HandleFunc("/npsanalysis/{userid}", npsAnalysis).Methods("GET")
	return r
}

func weeklyFTAnalysis(w http.ResponseWriter, r *http.Request) {
	zap.S().Info("***** fetch weekly FT Analysis *****")

	weeklyAnalysis, err := amplitudeapi.WeeklyMessage()
	if err != nil {
		zap.S().Errorf("error while fetching weekly FT analysis, error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Use this to prety print json
	jsonResp, err := json.MarshalIndent(weeklyAnalysis, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	//TODO: If we format message in Analysis bot then send weeklyAnalysis struct as resp.
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonResp); err != nil {
		zap.S().Errorf("Error while responding over http. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: Response to be sent not print statements.
	/*fmt.Sprintln("Total User signups :", weeklyAnalysis.Total_User_Signups)
	fmt.Sprintln("User who verified emails :", weeklyAnalysis.User_Who_Verified_Emails)
	fmt.Sprintln("Prep-node details")
	fmt.Sprintln("	Prep-node attempts")
	fmt.Sprintln("		New users:", weeklyAnalysis.PrepNode_Details.PrepNode_Attempts.New_Users)
	fmt.Sprintln("		Existing users:", weeklyAnalysis.PrepNode_Details.PrepNode_Attempts.Existing_Users)
	fmt.Sprintln("	Prep-node successes")
	fmt.Sprintln("		New users:", weeklyAnalysis.PrepNode_Details.PrepNode_Success.New_Users)
	fmt.Sprintln("		Existing users:", weeklyAnalysis.PrepNode_Details.PrepNode_Success.New_Users)
	fmt.Sprintln("	Prep-node errors")
	fmt.Sprintln("		New users:", weeklyAnalysis.PrepNode_Details.PrepNode_Errors.New_Users)
	fmt.Sprintln("		Existing users:", weeklyAnalysis.PrepNode_Details.PrepNode_Errors.Existing_Users)
	fmt.Sprintln("Cluster creation attempts")
	fmt.Sprintln("	New users:", weeklyAnalysis.Cluster_Creation_Attempts.New_Users)
	fmt.Sprintln("	Existing users:", weeklyAnalysis.Cluster_Creation_Attempts.Existing_Users)
	*/
}

func npsAnalysis(w http.ResponseWriter, r *http.Request) {
	zap.S().Info("***** NPS Analysis of user *****")

	vars := mux.Vars(r)
	userid := vars["userid"]
	zap.S().Debugf("UserID: %s", userid)
	amplitudeapi.NpsScoreAnalysis(userid)
	//TODO: Should merge both nps analysis and bugsnag analysis
	//bugsangapi.GetAllErrors()
}
