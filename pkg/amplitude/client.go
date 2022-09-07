package amplitude

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/platform9/ft-analyser-bot/pkg/config"
)

// Type definition for struct encapsulating amplitude APIs.
type api struct {
	client  *http.Client
	baseURL string
}

// getInfoAPI fetches the amplitude and charts info through http query.
func (api *api) getInfoAPI() ([]byte, error) {
	req, err := http.NewRequest("GET", api.baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Http request failed with error: %v", err)
	}

	req.Header.Add("Accept", "application/json")
	//amplitudeCreds := config.AmplitudeCreds()
	req.SetBasicAuth("userName", "PassWord")

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed with error: %v", err)
	}
	// TODO: Should handle the status codes i.e if not OK

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read the data, error: %v", err)
	}
	return data, nil
}

// getBorkInfoAPI fetches info using bork apis.
func (api *api) getBorkInfoAPI() ([]byte, error) {
	req, err := http.NewRequest("GET", api.baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Http request failed with error: %v", err)
	}
	// Bork needs to add token to the header for authorization of the request.
	bork_token_str := fmt.Sprintf("Bearer %s", config.BorkCreds())
	req.Header.Add("Authorization", bork_token_str)

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed with error: %v", err)
	}
	// TODO: Should handle the status codes i.e if not OK

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read the data, error: %v", err)
	}
	return data, nil
}
