package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/platform9/ft-analyser-bot/pkg/amplitude"
	"github.com/platform9/ft-analyser-bot/pkg/api"
	"github.com/platform9/ft-analyser-bot/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	version = "ft-analyser-bot version: v1.0"
	id      string
	email   string
)

func run(*cobra.Command, []string) {
	zap.S().Info("Starting FT analyser bot service...")
	zap.S().Infof("Version of FT analyser bot service: %s", version)
	fmt.Println("Analyzing user please wait...")
	npsAnalysis, err := amplitude.NpsScoreAnalysis(id, email)
	if err != nil {
		fmt.Println("Error", err)
	} else {
		out := api.GenNPSOutput(npsAnalysis)
		fmt.Println(out)
	}
	/*router := api.New()

	srv := &http.Server{
		Handler: router,
		Addr:    ":2112",
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			zap.S().Fatalf(err.Error())
		}
	}()

	go func() {
		ftBot.FtBotRun()
	}()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	select {
	case <-stop:
		zap.S().Info("server stopping...")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			zap.S().Fatalf(err.Error())
		}
	}*/
}

func main() {
	cmd := buildCmds()
	cmd.Execute()
	/*fmt.Println("NPS Analysis")
	fmt.Print("Enter userID: ")
	fmt.Scanf("%s", &id)
	fmt.Print("Enter emailID: ")
	fmt.Scanf("%s", &email)
	npsAnalysis, err := amplitude.NpsScoreAnalysis(id, email)
	if err != nil {
		fmt.Println("Error", err)
	} else {
		out := api.GenNPSOutput(npsAnalysis)
		fmt.Println(out)
	}*/

}

// Config file to read secrets like Amplitude, Hubspot, Bugsnag credentials.
var (
	homeDir, _ = os.UserHomeDir()
	// Change it accordingly
	analyserDir = filepath.Join(homeDir, "/envs")
	cfgFile     = filepath.Join(analyserDir, "config.yaml")
)

func buildCmds() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "analyze-nps",
		Short: "ft-analyser-bot is a service to generate FT weekly and NPS analysis",
		Long:  "ft-analyser-bot is a service to generate FT weekly and NPS analysis",
		Run:   run,
	}

	rootCmd.Flags().StringVar(&id, "id", "", "userID of user")
	rootCmd.Flags().StringVar(&email, "email", "", "emailID of user")
	return rootCmd
}

func initCfg() {
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		zap.S().Errorf(err.Error())
		panic(err)
	}
}

func init() {
	cobra.OnInitialize(initCfg)
	err := log.Logger()
	if err != nil {
		zap.S().Error("failed to initiate logger, Error is: %v", err.Error())
	}
}
