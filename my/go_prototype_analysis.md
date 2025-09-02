# TVerRec Go移植 プロトタイプ分析結果

## 現状の問題と解決策

### 🚨 発生している問題

1. **ネットワーク接続問題**
   - TVerのAPIへのアクセスがタイムアウト
   - yt-dlpコマンドも同様にタイムアウト
   - 地域制限またはネットワーク設定の問題の可能性

2. **作品情報取得の状況**
   - 自前のAPI呼び出し: ❌ タイムアウト
   - yt-dlpによる情報取得: ❌ タイムアウト（未確認）

### 💡 yt-dlpを活用したアプローチ

#### yt-dlpの情報取得機能
```bash
# 動画情報のみ取得（ダウンロードなし）
yt-dlp --dump-json --no-download <URL>

# メタデータ取得
yt-dlp --write-info-json --skip-download <URL>

# 利用可能フォーマット確認
yt-dlp --list-formats <URL>
```

#### 推奨アプローチの変更

**従来の設計（問題あり）:**
```
Go → TVerAPI直接呼び出し → 情報取得 → yt-dlp呼び出し → ダウンロード
```

**新しい設計（推奨）:**
```
Go → yt-dlp呼び出し → 情報取得+ダウンロード
```

### 🔄 修正されたプロトタイプ設計

#### 1. yt-dlp中心のアプローチ
```go
// 情報取得
func GetVideoInfoWithYtdlp(url string) (*VideoInfo, error) {
    cmd := exec.Command("yt-dlp", "--dump-json", "--no-download", url)
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    var info YtdlpInfo
    json.Unmarshal(output, &info)
    
    return convertToVideoInfo(info), nil
}

// ダウンロード
func DownloadVideo(url, outputDir string) error {
    cmd := exec.Command("yt-dlp", 
        "--output", filepath.Join(outputDir, "%(title)s.%(ext)s"),
        "--format", "best",
        url)
    return cmd.Run()
}
```

#### 2. yt-dlpのJSON出力構造
```json
{
  "id": "epuk32qiqy",
  "title": "エピソードタイトル",
  "uploader": "放送局名",
  "upload_date": "20231201",
  "description": "番組説明",
  "duration": 1440,
  "series": "シリーズ名",
  "season": "シーズン名",
  "episode": "エピソード名",
  "episode_number": 1
}
```

### 🛠 実装方針の変更

#### Phase 1: yt-dlp依存アプローチ
- ✅ **利点**: 
  - TVerの仕様変更に強い
  - 地域制限回避機能内蔵
  - 豊富なメタデータ取得
  - 安定したダウンロード機能

- ⚠️ **制約**:
  - yt-dlpへの依存
  - 情報取得の柔軟性が限定的

#### Phase 2: ハイブリッドアプローチ（将来）
- 基本: yt-dlp使用
- 拡張: 必要に応じて直接API呼び出し
- フォールバック: yt-dlpが失敗した場合の代替手段

### 📋 修正されたプロトタイプ仕様

#### 機能要件
1. **URL解析**: TVerのURL形式を解析
2. **情報取得**: yt-dlpの`--dump-json`で番組情報取得
3. **ダウンロード**: yt-dlpでの動画ダウンロード
4. **ファイル管理**: 出力ファイルの整理・リネーム
5. **エラーハンドリング**: yt-dlpのエラー処理

#### 技術スタック
- **情報取得**: yt-dlp --dump-json
- **ダウンロード**: yt-dlp 
- **プロセス管理**: os/exec
- **JSON処理**: encoding/json
- **ファイル操作**: path/filepath

### 🎯 次期プロトタイプの実装計画

#### 1. 基本機能
```go
type TVerDownloader struct {
    YtdlpPath   string
    OutputDir   string
    Options     []string
}

func (d *TVerDownloader) GetInfo(url string) (*VideoInfo, error)
func (d *TVerDownloader) Download(url string) error
func (d *TVerDownloader) DownloadWithInfo(url string) (*VideoInfo, error)
```

#### 2. 設定管理
```yaml
ytdlp:
  path: "yt-dlp"
  options:
    - "--format"
    - "best"
    - "--write-info-json"
  output_template: "%(series)s - %(episode)s - %(uploader)s.%(ext)s"
```

#### 3. エラーハンドリング
- yt-dlpプロセスの終了コード確認
- 標準エラー出力の解析
- リトライ機能
- フォールバック処理

### 🔍 検証すべき項目

1. **yt-dlpのTVer対応状況**
   - 対応エクストラクター確認
   - 取得可能メタデータの範囲
   - 地域制限対応

2. **パフォーマンス**
   - 情報取得速度
   - ダウンロード速度
   - メモリ使用量

3. **安定性**
   - エラー率
   - 長時間実行での安定性
   - 異常終了時の復旧

### 💭 結論

**現在の状況:**
- 直接API呼び出しは環境的な制約で困難
- yt-dlpを中心としたアプローチが現実的

**推奨方針:**
1. yt-dlpベースのプロトタイプ作成
2. 情報取得とダウンロードの統合
3. PowerShell版との機能比較
4. 段階的な機能拡張

この方針により、より実用的で安定したGo移植版を開発できると考えられます。