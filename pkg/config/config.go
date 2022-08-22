package config

import "github.com/spf13/viper"

type amplitudeCreds struct {
	ApiKey    string
	SecretKey string
}

type bugsnagCreds struct {
	ProjectID string
	AuthToken string
}

// Returns amplitude credentials
func AmplitudeCreds() amplitudeCreds {
	apiKey := viper.GetString("amplitude.apikey")
	secretKey := viper.GetString("amplitude.secretkey")
	return amplitudeCreds{ApiKey: apiKey, SecretKey: secretKey}
}

func BorkCreds() string {
	return viper.GetString("bork.token")
}

func BugsnagCreds() bugsnagCreds {
	projectID := viper.GetString("bugsnag.projectID")
	authToken := viper.GetString("bugsnag.authToken")
	return bugsnagCreds{ProjectID: projectID, AuthToken: authToken}
}
