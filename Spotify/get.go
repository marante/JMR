package Spotify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Get gets the resources for the given URL string, only if token is valid.
func Get(url string, token string, result interface{}) error {
	for {
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating a new request", err)
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == rateLimitExceededStatusCode {
			time.Sleep(retryDuration(resp))
			continue
		}

		if resp.StatusCode != http.StatusOK {
			return decodeError(resp)
		}

		err = json.NewDecoder(resp.Body).Decode(result)
		if err != nil {
			return err
		}

		break
	}
	return nil
}
