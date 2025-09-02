package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// TVerClient manages communication with the TVer API.
type TVerClient struct {
	PlatformUID   string
	PlatformToken string
	HTTPClient    *http.Client
}

// NewTVerClient creates a new TVer API client.
func NewTVerClient() *TVerClient {
	return &TVerClient{
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// GetToken fetches an authentication token from the TVer platform API.
func (c *TVerClient) GetToken() error {
	url := "https://platform-api.tver.jp/v2/api/platform_users/browser/create"

	req, err := http.NewRequest("POST", url, strings.NewReader("device_type=pc"))
	if err != nil {
		return fmt.Errorf("リクエスト作成エラー: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("APIリクエストエラー: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("APIエラー: ステータスコード %d, レスポンス: %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		Result struct {
			PlatformUID   string `json:"platform_uid"`
			PlatformToken string `json:"platform_token"`
		} `json:"Result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("レスポンス解析エラー: %w", err)
	}

	c.PlatformUID = tokenResp.Result.PlatformUID
	c.PlatformToken = tokenResp.Result.PlatformToken

	fmt.Printf("トークン取得完了: UID=%s\n", c.PlatformUID[:8]+"...")
	return nil
}

// GetSeriesSeasons fetches a list of season IDs for a given series ID.
func (c *TVerClient) GetSeriesSeasons(seriesID string) ([]string, error) {
	url := fmt.Sprintf("https://platform-api.tver.jp/service/api/v1/callSeriesSeasons/%s?platform_uid=%s&platform_token=%s",
		seriesID, c.PlatformUID, c.PlatformToken)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成エラー: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("APIリクエストエラー: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("APIエラー: ステータスコード %d, レスポンス: %s", resp.StatusCode, string(body))
	}

	var apiResp struct {
		Result struct {
			Contents []struct {
				Type    string `json:"Type"`
				Content struct {
					ID string `json:"Id"`
				} `json:"Content"`
			} `json:"Contents"`
		} `json:"Result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("レスポンス解析エラー: %w", err)
	}

	var seasonIDs []string
	for _, content := range apiResp.Result.Contents {
		if content.Type == "season" {
			seasonIDs = append(seasonIDs, content.Content.ID)
		}
	}

	return seasonIDs, nil
}

// GetSeasonEpisodes fetches a list of episodes for a given season ID.
func (c *TVerClient) GetSeasonEpisodes(seasonID string) ([]EpisodeEntry, error) {
	url := fmt.Sprintf("https://platform-api.tver.jp/service/api/v1/callSeasonEpisodes/%s?platform_uid=%s&platform_token=%s",
		seasonID, c.PlatformUID, c.PlatformToken)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成エラー: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("APIリクエストエラー: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("APIエラー: ステータスコード %d, レスポンス: %s", resp.StatusCode, string(body))
	}

	var apiResp struct {
		Result struct {
			Contents []struct {
				Type    string `json:"Type"`
				Content struct {
					ID    string `json:"Id"`
					Title string `json:"Title"`
					EndAt int64  `json:"EndAt"`
				} `json:"Content"`
			} `json:"Contents"`
		} `json:"Result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("レスポンス解析エラー: %w", err)
	}

	var episodes []EpisodeEntry
	for _, content := range apiResp.Result.Contents {
		if content.Type == "episode" {
			episode := EpisodeEntry{
				Type:       "video",
				Title:      content.Content.Title,
				WebpageURL: fmt.Sprintf("https://tver.jp/episodes/%s", content.Content.ID),
				ID:         content.Content.ID,
				Extractor:  "TVer",
			}
			episodes = append(episodes, episode)
		}
	}

	return episodes, nil
}
