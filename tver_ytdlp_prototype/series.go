// series.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
)

// min関数（Go 1.21未満の場合）
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// シリーズ情報（yt-dlpから取得）
type SeriesInfo struct {
	Type      string         `json:"_type"`
	Entries   []EpisodeEntry `json:"entries"`
	Title     string         `json:"title"`
	ID        string         `json:"id"`
	Extractor string         `json:"extractor"`
}

// エピソード情報
type EpisodeEntry struct {
	Type       string `json:"_type"`
	Title      string `json:"title"`
	WebpageURL string `json:"webpage_url"`
	ID         string `json:"id"`
	Extractor  string `json:"extractor"`
}

// 解析済みエピソード情報
type ParsedEpisode struct {
	EpisodeNumber int
	Title         string
	URL           string
	ID            string
	OriginalTitle string
}

// シリーズ管理
type SeriesManager struct {
	YtdlpPath string
}

// 新しいシリーズマネージャーを作成
func NewSeriesManager() *SeriesManager {
	return &SeriesManager{
		YtdlpPath: "yt-dlp",
	}
}

// シリーズURLからエピソード一覧を取得（TVerAPI使用）
func (sm *SeriesManager) GetSeriesInfo(seriesURL string) (*SeriesInfo, error) {
	fmt.Printf("シリーズ情報取得開始: %s\n", seriesURL)

	seriesID, err := sm.extractSeriesID(seriesURL)
	if err != nil {
		return nil, fmt.Errorf("シリーズID抽出エラー: %w", err)
	}
	fmt.Printf("シリーズID: %s\n", seriesID)

	client := NewTVerClient()
	if err := client.GetToken(); err != nil {
		return nil, fmt.Errorf("トークン取得エラー: %w", err)
	}

	seasons, err := client.GetSeriesSeasons(seriesID)
	if err != nil {
		return nil, fmt.Errorf("シーズン取得エラー: %w", err)
	}
	fmt.Printf("シーズン数: %d\n", len(seasons))

	var allEpisodes []EpisodeEntry
	for _, seasonID := range seasons {
		episodes, err := client.GetSeasonEpisodes(seasonID)
		if err != nil {
			fmt.Printf("シーズン %s のエピソード取得エラー: %v\n", seasonID, err)
			continue
		}
		allEpisodes = append(allEpisodes, episodes...)
	}

	seriesInfo := &SeriesInfo{
		Type:    "playlist",
		Entries: allEpisodes,
		Title:   "TVerシリーズ",
	}

	fmt.Printf("エピソード数: %d話\n", len(allEpisodes))
	return seriesInfo, nil
}

// シリーズURLからシリーズIDを抽出
func (sm *SeriesManager) extractSeriesID(seriesURL string) (string, error) {
	re := regexp.MustCompile(`series/([a-zA-Z0-9]+)`)
	matches := re.FindStringSubmatch(seriesURL)
	if len(matches) < 2 {
		return "", fmt.Errorf("URLからシリーズIDを抽出できません: %s", seriesURL)
	}
	return matches[1], nil
}

// タイトルからエピソード番号を抽出
func extractEpisodeNumber(title string) int {
	// パターン1: "第X話"
	re1 := regexp.MustCompile(`第(\d+)話`)
	if matches := re1.FindStringSubmatch(title); len(matches) >= 2 {
		if num, err := strconv.Atoi(matches[1]); err == nil {
			return num
		}
	}

	// パターン2: "Episode X"
	re2 := regexp.MustCompile(`Episode\s+(\d+)`)
	if matches := re2.FindStringSubmatch(title); len(matches) >= 2 {
		if num, err := strconv.Atoi(matches[1]); err == nil {
			return num
		}
	}

	// パターン3: "#X"
	re3 := regexp.MustCompile(`#(\d+)`)
	if matches := re3.FindStringSubmatch(title); len(matches) >= 2 {
		if num, err := strconv.Atoi(matches[1]); err == nil {
			return num
		}
	}

	return 0 // 抽出できない場合
}

// エピソード一覧を解析・整理
func (sm *SeriesManager) ParseEpisodes(seriesInfo *SeriesInfo) []ParsedEpisode {
	var episodes []ParsedEpisode

	for _, entry := range seriesInfo.Entries {
		episodeNum := extractEpisodeNumber(entry.Title)

		// URLからエピソードIDを抽出
		episodeID := ""
		if re := regexp.MustCompile(`episodes/([a-zA-Z0-9]+)`); re != nil {
			if matches := re.FindStringSubmatch(entry.WebpageURL); len(matches) >= 2 {
				episodeID = matches[1]
			}
		}

		episode := ParsedEpisode{
			EpisodeNumber: episodeNum,
			Title:         entry.Title,
			URL:           entry.WebpageURL,
			ID:            episodeID,
			OriginalTitle: entry.Title,
		}

		episodes = append(episodes, episode)
	}

	// エピソード番号で昇順ソート
	sort.Slice(episodes, func(i, j int) bool {
		// エピソード番号が0の場合は最後に
		if episodes[i].EpisodeNumber == 0 && episodes[j].EpisodeNumber != 0 {
			return false
		}
		if episodes[i].EpisodeNumber != 0 && episodes[j].EpisodeNumber == 0 {
			return true
		}
		return episodes[i].EpisodeNumber < episodes[j].EpisodeNumber
	})

	return episodes
}

// 指定範囲のエピソードをフィルタリング
func (sm *SeriesManager) FilterEpisodes(episodes []ParsedEpisode, fromEpisode, toEpisode int) []ParsedEpisode {
	var filtered []ParsedEpisode

	for _, ep := range episodes {
		// エピソード番号が抽出できない場合はスキップ
		if ep.EpisodeNumber == 0 {
			continue
		}

		// 範囲チェック
		if fromEpisode > 0 && ep.EpisodeNumber < fromEpisode {
			continue
		}
		if toEpisode > 0 && ep.EpisodeNumber > toEpisode {
			continue
		}

		filtered = append(filtered, ep)
	}

	return filtered
}

// エピソード一覧を表示
func (sm *SeriesManager) DisplayEpisodes(episodes []ParsedEpisode) {
	fmt.Println("\n=== エピソード一覧 ===")
	for i, ep := range episodes {
		if ep.EpisodeNumber > 0 {
			fmt.Printf("%2d. 第%d話: %s\n", i+1, ep.EpisodeNumber, ep.Title)
		} else {
			fmt.Printf("%2d. [番号不明]: %s\n", i+1, ep.Title)
		}
		fmt.Printf("    ID: %s\n", ep.ID)
		fmt.Printf("    URL: %s\n", ep.URL)
		fmt.Println()
	}
	fmt.Printf("合計: %d話\n", len(episodes))
	fmt.Println("==================")
}

// シリーズ情報をJSONファイルに保存
func (sm *SeriesManager) SaveSeriesToFile(episodes []ParsedEpisode, filename string) error {
	data := map[string]interface{}{
		"episodes": episodes,
		"count":    len(episodes),
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("ファイル作成エラー: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("JSON書き込みエラー: %w", err)
	}

	fmt.Printf("シリーズ情報を保存: %s\n", filename)
	return nil
}
