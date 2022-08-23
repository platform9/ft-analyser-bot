package api

import (
	"fmt"

	amplitudeapi "github.com/platform9/ft-analyser-bot/pkg/amplitude"
)

func GenerateOutputString(weeklyAnalysis amplitudeapi.WeeklyAnalysis) string {
	var s string
	s += fmt.Sprintln("Total User signups :", weeklyAnalysis.Total_User_Signups)
	s += fmt.Sprintln("User who verified emails :", weeklyAnalysis.User_Who_Verified_Emails)
	s += fmt.Sprintln("Prep-node details")
	s += fmt.Sprintln("	- Prep-node attempts")
	s += fmt.Sprintln("		- New users:", weeklyAnalysis.PrepNode_Details.PrepNode_Attempts.New_Users)
	s += fmt.Sprintln("		- Existing users:", weeklyAnalysis.PrepNode_Details.PrepNode_Attempts.Existing_Users)
	s += fmt.Sprintln("	- Prep-node successes")
	s += fmt.Sprintln("		- New users:", weeklyAnalysis.PrepNode_Details.PrepNode_Success.New_Users)
	s += fmt.Sprintln("		- Existing users:", weeklyAnalysis.PrepNode_Details.PrepNode_Success.New_Users)
	s += fmt.Sprintln("	- Prep-node errors")
	s += fmt.Sprintln("		- New users:", weeklyAnalysis.PrepNode_Details.PrepNode_Errors.New_Users)
	s += fmt.Sprintln("		- Existing users:", weeklyAnalysis.PrepNode_Details.PrepNode_Errors.Existing_Users)
	s += fmt.Sprintln("Cluster creation attempts")
	s += fmt.Sprintln("	- New users:", weeklyAnalysis.Cluster_Creation_Attempts.New_Users)
	s += fmt.Sprintln("	- Existing users:", weeklyAnalysis.Cluster_Creation_Attempts.Existing_Users)
	return s
}
