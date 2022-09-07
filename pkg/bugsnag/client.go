package bugsnag

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Type definition for struct encapsulating amplitude APIs.
type api struct {
	client  *http.Client
	baseURL string
}

// getBugsnagInfoAPI fetches info using bork apis.
func (api *api) getBugsnagInfoAPI() ([]byte, error) {
	var token = "token <TOKEN>"
	req, err := http.NewRequest("GET", api.baseURL, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, fmt.Errorf("Http request failed with error: %v", err)
	}

	req.Header.Add("X-Version", "2")
	req.Header.Add("Authorization", token)
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
