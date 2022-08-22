package main

import (
	"fmt"

	amplitudeapi "github.com/platform9/ft-analyser-bot/amplitudeAPI"
)

func main() {
	fmt.Println("Analysis")
	userID := "530388e91baf49f49c9cd059764fd7b3"
	//amplitudeapi.WeeklyMessage()
	amplitudeapi.NPS_Score_Analysis(userID)
}
