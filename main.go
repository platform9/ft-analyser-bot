package main

import (
	"fmt"

	amplitudeapi "github.com/platform9/ft-analyser-bot/amplitudeAPI"
)

func main() {
	fmt.Println("Analysis")
	userID := "1529858a0f324051a71878724eac1f82"
	//amplitudeapi.WeeklyMessage()
	amplitudeapi.NPS_Score_Analysis(userID)
}
