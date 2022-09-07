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
		zap.S().Errorf("Error while marshalling the response: %v", err)
		return
	}

	//out := GenerateOutputString(weeklyAnalysis)
	//TODO: If we format message in Analysis bot then send weeklyAnalysis struct as resp.
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonResp); err != nil {
		zap.S().Errorf("Error while responding over http. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func npsAnalysis(w http.ResponseWriter, r *http.Request) {
	zap.S().Info("***** NPS Analysis of user *****")

	vars := mux.Vars(r)
	userid := vars["userid"]
	zap.S().Debugf("UserID: %s", userid)
	npsAnalysis, err := amplitudeapi.NpsScoreAnalysis(userid, "")
	if err != nil {
		zap.S().Errorf("error while fetching nps analysis for user, error: %v", userid, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Use this to prety print json
	jsonResp, err := json.MarshalIndent(npsAnalysis, "", "  ")
	if err != nil {
		fmt.Println(err)
		zap.S().Errorf("Error while marshalling the response: %v", err)
		return
	}
	//out := GenNPSOutput(npsAnalysis)

	//TODO: If we format message in Analysis bot then send weeklyAnalysis struct as resp.
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonResp); err != nil {
		zap.S().Errorf("Error while responding over http. Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
