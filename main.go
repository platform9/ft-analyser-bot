package main

import (
	"fmt"

	bugsangapi "github.com/platform9/ft-analyser-bot/bugsnag"
)

func main() {
	fmt.Println("Analysis")

	//amplitudeapi.WeeklyMessage()
	//amplitudeapi.NPS_Score_Analysis("userID")
	bugsangapi.GetAllErrors("userID")
}
