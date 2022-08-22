package main

import (
	"fmt"

	bugsangapi "github.com/platform9/ftanalysis/bugsnag"
)

func main() {
	fmt.Println("Analysis")

	//amplitudeapi.WeeklyMessage()
	//amplitudeapi.NPS_Score_Analysis("5338d68826cf4082908a2aea5758ee3d")
	bugsangapi.GetAllErrors()
}
