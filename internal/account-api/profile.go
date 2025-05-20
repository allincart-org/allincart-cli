package account_api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/allincart-org/allincart-cli/logging"
)

func (c *Client) GetMyProfile(ctx context.Context) (*MyProfile, error) {
	errorFormat := "GetMyProfile: %v"

	request, err := c.NewAuthenticatedRequest(ctx, "GET", fmt.Sprintf("%s/account/profile", ApiUrl), nil)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			logging.FromContext(ctx).Errorf("GetMyProfile: %v", err)
		}
	}()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(errorFormat, err)
	}

	var profile MyProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	return &profile, nil
}

type MyProfile struct {
	PersonalData struct {
		Name string `json:"name"`
	} `json:"personalData"`
}
