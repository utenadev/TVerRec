package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// yt-dlpから取得する動画情報
type YtdlpVideoInfo struct {
	ID            string  `json:"id"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	Uploader      string  `json:"uploader"`
	UploaderID    string  `json:"uploader_id"`
	UploadDate    string  `json:"upload_date"`
	Duration      float64 `json:"duration"`
	Series        string  `json:"series"`
	Season        string  `json:"season"`
	Episode       string  `json:"episode"`
	EpisodeNumber int     `json:"episode_number"`
	Webpage       string  `json:"webpage_url"`
	Extractor     string  `json:"extractor"`
	ExtractorKey  string  `json:"extractor_key"`
}

// TVerダウンローダー
type TVerDownloader struct {
	YtdlpPath string
	OutputDir string
	Options   []string
}

// 新しいTVerダウンローダーを作成
func NewTVerDownloader(outputDir string) *TVerDownloader {
	return &TVerDownloader{
		YtdlpPath: "yt-dlp",
		OutputDir: outputDir,
		Options: []string{
			"-N", "10", // 10並列ダウンロード
			"--write-info-json", // 情報JSONファイルも出力
		},
	}
}

// URLからエピソードIDを抽出
func extractEpisodeID(url string) (string, error) {
	re := regexp.MustCompile(`episodes/([a-zA-Z0-9]+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", fmt.Errorf("URLからエピソードIDを抽出できません: %s", url)
	}
	return matches[1], nil
}

// yt-dlpを使って動画情報のみを取得
func (d *TVerDownloader) GetVideoInfo(url string) (*YtdlpVideoInfo, error) {
	fmt.Printf("動画情報取得開始: %s\n", url)
	
	// yt-dlpコマンドを構築（情報取得のみ）
	args := []string{
		"--dump-json",
		"--no-download",
		url,
	}
	
	cmd := exec.Command(d.YtdlpPath, args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("yt-dlp情報取得エラー: %w", err)
	}
	
	var info YtdlpVideoInfo
	if err := json.Unmarshal(output, &info); err != nil {
		return nil, fmt.Errorf("JSON解析エラー: %w", err)
	}
	
	fmt.Printf("情報取得完了: %s\n", info.Title)
	return &info, nil
}

// yt-dlpを使って動画をダウンロード
func (d *TVerDownloader) DownloadVideo(url string) error {
	fmt.Printf("ダウンロード開始: %s\n", url)
	
	// 出力テンプレートを設定
	outputTemplate := filepath.Join(d.OutputDir, "%(series)s - %(episode)s - %(uploader)s.%(ext)s")
	
	// yt-dlpコマンドを構築
	args := append(d.Options,
		"-o", outputTemplate,
		url,
	)
	
	cmd := exec.Command(d.YtdlpPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	start := time.Now()
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("yt-dlpダウンロードエラー: %w", err)
	}
	
	duration := time.Since(start)
	fmt.Printf("ダウンロード完了 (所要時間: %v)\n", duration)
	return nil
}

// 動画情報とダウンロードを同時実行
func (d *TVerDownloader) GetInfoAndDownload(url string) (*YtdlpVideoInfo, error) {
	// まず情報を取得
	info, err := d.GetVideoInfo(url)
	if err != nil {
		return nil, fmt.Errorf("情報取得失敗: %w", err)
	}
	
	// 情報を表示
	d.displayVideoInfo(info)
	
	// ダウンロード実行
	if err := d.DownloadVideo(url); err != nil {
		return info, fmt.Errorf("ダウンロード失敗: %w", err)
	}
	
	return info, nil
}

// 動画情報を表示
func (d *TVerDownloader) displayVideoInfo(info *YtdlpVideoInfo) {
	fmt.Println("\n=== 動画情報 ===")
	fmt.Printf("ID: %s\n", info.ID)
	fmt.Printf("タイトル: %s\n", info.Title)
	fmt.Printf("シリーズ: %s\n", info.Series)
	fmt.Printf("シーズン: %s\n", info.Season)
	fmt.Printf("エピソード: %s\n", info.Episode)
	if info.EpisodeNumber > 0 {
		fmt.Printf("エピソード番号: %d\n", info.EpisodeNumber)
	}
	fmt.Printf("配信者: %s\n", info.Uploader)
	fmt.Printf("配信日: %s\n", info.UploadDate)
	if info.Duration > 0 {
		fmt.Printf("長さ: %.0f秒 (%.1f分)\n", info.Duration, info.Duration/60)
	}
	fmt.Printf("URL: %s\n", info.Webpage)
	if info.Description != "" {
		fmt.Printf("説明: %s\n", strings.TrimSpace(info.Description))
	}
	fmt.Println("================\n")
}

// 情報をJSONファイルに保存
func (d *TVerDownloader) SaveInfoToFile(info *YtdlpVideoInfo) error {
	filename := fmt.Sprintf("%s_info.json", info.ID)
	filepath := filepath.Join(d.OutputDir, filename)
	
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("ファイル作成エラー: %w", err)
	}
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(info); err != nil {
		return fmt.Errorf("JSON書き込みエラー: %w", err)
	}
	
	fmt.Printf("動画情報を保存: %s\n", filepath)
	return nil
}

// yt-dlpの存在確認
func checkYtdlp() error {
	cmd := exec.Command("yt-dlp", "--version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("yt-dlpが見つかりません。インストールしてください: %w", err)
	}
	
	version := strings.TrimSpace(string(output))
	fmt.Printf("yt-dlp バージョン: %s\n", version)
	return nil
}

// 使用方法を表示
func showUsage() {
	fmt.Println("TVerアニメダウンローダー (yt-dlpベース)")
	fmt.Println()
	fmt.Println("使用方法:")
	fmt.Println("  go run *.go <command> <TVerURL> [オプション] [出力ディレクトリ]")
	fmt.Println()
	fmt.Println("コマンド:")
	fmt.Println("  info     - 動画情報のみ取得")
	fmt.Println("  download - 動画をダウンロード")
	fmt.Println("  both     - 情報取得とダウンロードの両方")
	fmt.Println("  series   - シリーズ情報取得・一括ダウンロード")
	fmt.Println()
	fmt.Println("シリーズオプション:")
	fmt.Println("  --list           - エピソード一覧のみ表示")
	fmt.Println("  --from N         - N話以降をダウンロード")
	fmt.Println("  --to N           - N話まででダウンロード")
	fmt.Println("  --all            - 全話ダウンロード")
	fmt.Println()
	fmt.Println("例:")
	fmt.Println("  go run *.go info https://tver.jp/episodes/epuk32qiqy")
	fmt.Println("  go run *.go series https://tver.jp/series/srrazrs5j2 --list")
	fmt.Println("  go run *.go series https://tver.jp/series/srrazrs5j2 --from 10")
	fmt.Println("  go run *.go series https://tver.jp/series/srrazrs5j2 --from 10 --to 15")
}

func main() {
	if len(os.Args) < 3 {
		showUsage()
		os.Exit(1)
	}
	
	command := os.Args[1]
	targetURL := os.Args[2]
	
	// オプション解析
	var fromEpisode, toEpisode int
	var listOnly, allEpisodes bool
	outputDir := "./downloads"
	
	for i := 3; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch {
		case arg == "--list":
			listOnly = true
		case arg == "--all":
			allEpisodes = true
		case arg == "--from" && i+1 < len(os.Args):
			if num, err := strconv.Atoi(os.Args[i+1]); err == nil {
				fromEpisode = num
				i++ // 次の引数をスキップ
			}
		case arg == "--to" && i+1 < len(os.Args):
			if num, err := strconv.Atoi(os.Args[i+1]); err == nil {
				toEpisode = num
				i++ // 次の引数をスキップ
			}
		case !strings.HasPrefix(arg, "--"):
			outputDir = arg
		}
	}
	
	// yt-dlpの存在確認
	if err := checkYtdlp(); err != nil {
		log.Fatalf("yt-dlp確認エラー: %v", err)
	}
	
	// 出力ディレクトリを作成
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("出力ディレクトリ作成エラー: %v", err)
	}
	
	fmt.Printf("出力ディレクトリ: %s\n", outputDir)
	fmt.Println()
	
	// コマンドに応じて処理を実行
	switch command {
	case "info":
		// エピソードIDを抽出
		episodeID, err := extractEpisodeID(targetURL)
		if err != nil {
			log.Fatalf("エピソードID抽出エラー: %v", err)
		}
		fmt.Printf("エピソードID: %s\n", episodeID)
		
		// ダウンローダーを初期化
		downloader := NewTVerDownloader(outputDir)
		info, err := downloader.GetVideoInfo(targetURL)
		if err != nil {
			log.Fatalf("情報取得エラー: %v", err)
		}
		downloader.displayVideoInfo(info)
		if err := downloader.SaveInfoToFile(info); err != nil {
			log.Printf("情報保存エラー: %v", err)
		}
		
	case "download":
		// エピソードIDを抽出
		episodeID, err := extractEpisodeID(targetURL)
		if err != nil {
			log.Fatalf("エピソードID抽出エラー: %v", err)
		}
		fmt.Printf("エピソードID: %s\n", episodeID)
		
		// ダウンローダーを初期化
		downloader := NewTVerDownloader(outputDir)
		if err := downloader.DownloadVideo(targetURL); err != nil {
			log.Fatalf("ダウンロードエラー: %v", err)
		}
		
	case "both":
		// エピソードIDを抽出
		episodeID, err := extractEpisodeID(targetURL)
		if err != nil {
			log.Fatalf("エピソードID抽出エラー: %v", err)
		}
		fmt.Printf("エピソードID: %s\n", episodeID)
		
		// ダウンローダーを初期化
		downloader := NewTVerDownloader(outputDir)
		info, err := downloader.GetInfoAndDownload(targetURL)
		if err != nil {
			log.Fatalf("処理エラー: %v", err)
		}
		if err := downloader.SaveInfoToFile(info); err != nil {
			log.Printf("情報保存エラー: %v", err)
		}
		
	case "series":
		// シリーズマネージャーを初期化
		seriesManager := NewSeriesManager()
		
		// シリーズ情報を取得
		seriesInfo, err := seriesManager.GetSeriesInfo(targetURL)
		if err != nil {
			log.Fatalf("シリーズ情報取得エラー: %v", err)
		}
		
		// エピソードを解析
		episodes := seriesManager.ParseEpisodes(seriesInfo)
		
		// 範囲フィルタリング
		if !allEpisodes {
			episodes = seriesManager.FilterEpisodes(episodes, fromEpisode, toEpisode)
		}
		
		// エピソード一覧を表示
		seriesManager.DisplayEpisodes(episodes)
		
		// シリーズ情報をファイルに保存
		seriesFile := filepath.Join(outputDir, "series_info.json")
		if err := seriesManager.SaveSeriesToFile(episodes, seriesFile); err != nil {
			log.Printf("シリーズ情報保存エラー: %v", err)
		}
		
		// リスト表示のみの場合は終了
		if listOnly {
			fmt.Println("エピソード一覧表示完了!")
			return
		}
		
		// ダウンロード実行
		if len(episodes) == 0 {
			fmt.Println("ダウンロード対象のエピソードがありません。")
			return
		}
		
		fmt.Printf("\n%d話のダウンロードを開始します...\n", len(episodes))
		downloader := NewTVerDownloader(outputDir)
		
		for i, episode := range episodes {
			fmt.Printf("\n[%d/%d] ダウンロード中: %s\n", i+1, len(episodes), episode.Title)
			
			if err := downloader.DownloadVideo(episode.URL); err != nil {
				log.Printf("エピソード %d のダウンロードエラー: %v", episode.EpisodeNumber, err)
				continue
			}
			
			fmt.Printf("完了: %s\n", episode.Title)
		}
		
	default:
		fmt.Printf("不明なコマンド: %s\n", command)
		showUsage()
		os.Exit(1)
	}
	
	fmt.Println("処理完了!")
}