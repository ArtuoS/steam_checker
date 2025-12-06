package steam

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	. "steam_checker/internal/shared/schema/steam"
	"strconv"
)

type Integration struct {
	APIKey      string
	APIURL      string
	APIStoreURL string
}

func NewIntegration() *Integration {
	return &Integration{
		APIURL:      "https://api.steampowered.com",
		APIStoreURL: "https://store.steampowered.com",
		APIKey:      os.Getenv("STEAM_API_KEY"),
	}
}

func (i *Integration) GetPlayerCount(ctx context.Context, appID int) (GetPlayerCountData, error) {
	endpoint := fmt.Sprintf("%s/ISteamUserStats/GetNumberOfCurrentPlayers/v1/", i.APIURL)

	params := url.Values{}
	params.Set("key", i.APIKey)
	params.Set("appid", string(appID))
	params.Set("format", "json")

	reqURL := endpoint + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return GetPlayerCountData{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GetPlayerCountData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return GetPlayerCountData{}, fmt.Errorf("steam api returned status %d", resp.StatusCode)
	}

	var res GetPlayerCountResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return GetPlayerCountData{}, err
	}

	return res.Data, nil
}

func (i *Integration) GetAppDetails(ctx context.Context, appID int) (GetAppDetailsData, error) {
	endpoint := fmt.Sprintf("%s/api/appdetails/", i.APIStoreURL)

	params := url.Values{}
	params.Set("appids", strconv.Itoa(appID))

	reqURL := endpoint + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return GetAppDetailsData{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GetAppDetailsData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return GetAppDetailsData{}, fmt.Errorf("steam api returned status %d", resp.StatusCode)
	}

	var res GetAppDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return GetAppDetailsData{}, err
	}

	fmt.Println(res["Data"].Data)

	return res["Data"].Data, nil
}
