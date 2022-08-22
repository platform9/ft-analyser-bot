package main

import (
	"fmt"

	bugsangapi "github.com/platform9/ft-analyser-bot/bugsnag"
)

func main() {
	fmt.Println("Analysis")

	//amplitudeapi.WeeklyMessage()
	//amplitudeapi.NPS_Score_Analysis("5338d68826cf4082908a2aea5758ee3d")
	bugsangapi.GetAllErrors("id")
}
