package api

import (
	"fmt"
	"strings"

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
	s += fmt.Sprintln("		- Existing users:", weeklyAnalysis.PrepNode_Details.PrepNode_Success.Existing_Users)
	s += fmt.Sprintln("	- Prep-node errors")
	s += fmt.Sprintln("		- New users:", weeklyAnalysis.PrepNode_Details.PrepNode_Errors.New_Users)
	s += fmt.Sprintln("		- Existing users:", weeklyAnalysis.PrepNode_Details.PrepNode_Errors.Existing_Users)
	s += fmt.Sprintln("Cluster creation attempts")
	s += fmt.Sprintln("	- New users:", weeklyAnalysis.Cluster_Creation_Attempts.New_Users)
	s += fmt.Sprintln("	- Existing users:", weeklyAnalysis.Cluster_Creation_Attempts.Existing_Users)
	return s
}

func GenNPSOutput(npsAnalysis amplitudeapi.NPSAnalysis) string {
	var s string
	s += fmt.Sprintln("* User is from :", npsAnalysis.UserCountry)
	s += fmt.Sprintf("* Active from %s, last seen %s", npsAnalysis.FirstSeen, npsAnalysis.LastSeen)
	s += fmt.Sprintf("\n* DU %s is hosted on the %s cluster", npsAnalysis.HostDetails.FQDN, strings.Trim(npsAnalysis.HostDetails.BorkCluster, ".platform9.io"))
	s += fmt.Sprintf("\n* No of hosts attached %s, and active %s", npsAnalysis.HostDetails.HostCount, npsAnalysis.HostDetails.ActiveHosts)
	s += fmt.Sprintf("\n* User performed Check-Node successfully %d time/s", npsAnalysis.CLIEvents.ChecknodeSuccess)
	s += fmt.Sprintf("\n* User preformed Prep-Node successfully %d time/s, failed %d time/s", npsAnalysis.CLIEvents.Prepnode.Success, npsAnalysis.CLIEvents.Prepnode.Failure)

	if npsAnalysis.CLIEvents.Prepnode.Failure != 0 {
		s += fmt.Sprintln("\n	- Prepnode failed due to errors: ")
		for _, val := range npsAnalysis.CLIEvents.PrepnodeErrors {
			s += fmt.Sprintln("	 - ", val)
		}
	}

	s += fmt.Sprintf("\n* Created cluster successfully %d times, and deleted %d times\n", npsAnalysis.CLIEvents.ClusterCreation.Success, npsAnalysis.CLIEvents.ClusterCreation.Delete)
	s += fmt.Sprintln("* Few of UI activites are: ")
	for _, val := range npsAnalysis.UserActivities {
		s += fmt.Sprintln("	 - ", val)
	}

	s += fmt.Sprintln("* UI errors are: ")
	if len(npsAnalysis.UIErrors) == 0 {
		s += fmt.Sprintf("nil")
	} else {
		for _, val := range npsAnalysis.UIErrors {
			s += fmt.Sprintln("	 - ", val)
		}
	}

	return s
}
