# Go言語への移植検討・整理

## 1. 概要

このドキュメントは、PowerShellで実装されたTVerRecをGo言語に移植する際の検討内容と整理結果をまとめたものです。移植の背景、アーキテクチャ設計、作業タスクリストなどを記載しています。

## 2. 移植の背景と目的

- **保守性向上**: Goは静的型付け言語であり、コンパイル時チェックによりエラーを早期発見できる。
- **パフォーマンス**: コンパイル済みバイナリとして動作するため、実行速度が速い。
- **依存関係の簡素化**: 単一バイナリとして配布可能で、実行環境への依存を減らせる。
- **クロスプラットフォーム対応**: Windows/Mac/Linux向けに同一ソースコードからバイナリを生成可能。

## 3. TVerRecコード分析

### PowerShellスクリプト構成

- **`src/functions/`**: TVer APIクライアント、動画情報取得、ダウンロードなどコア機能を提供
- **`src/gui/`**: GUIアプリケーション用スクリプト（WPF/XAML）
- **`src/`**: 各種機能スクリプト（一括ダウンロード、リスト生成、バリデーションなど）
- **設定管理**: JSON/XMLファイル、ユーザー設定スクリプト
- **外部ツール依存**: `youtube-dl`または`yt-dlp`

### 主要機能

1. TVer APIトークン取得 (`Get-Token`)
2. キーワード/URLから番組情報取得 (`Get-VideoLinksFromKeyword`, `Get-VideoInfo`)
3. シリーズ情報取得とエピソード一覧解析 (`Get-SeriesInfo`)
4. 動画ダウンロード (`Start-DownloadProcess`)
5. 設定管理 (`Initialize-Configuration`, `Save-Settings`)
6. GUI操作 (`Show-GUI`, `Update-GUIProgress`)

## 4. TVer APIの挙動調査 (yt-dlp出力解析)

### 動画情報取得フロー

1. **APIトークン取得**
   - URL: `https://platform-api.tver.jp/v2/api/platform_users/browser/create`
   - Method: POST
   - Headers: 特定のUser-AgentとForwardedヘッダーが必要

2. **シリーズ情報取得**
   - URL: `https://platform-api.tver.jp/service/api/v1/callSeriesSeasons/{series_id}`
   - Query Params: `platform_uid`, `platform_token`
   - Response: シーズンIDリスト

3. **エピソード情報取得**
   - URL: `https://platform-api.tver.jp/service/api/v1/callSeasonEpisodes/{season_id}`
   - Query Params: `platform_uid`, `platform_token`
   - Response: エピソードIDリストとメタデータ

4. **動画詳細情報取得**
   - URL: `https://platform-api.tver.jp/service/api/v1/callEpisode/{episode_id}`
   - Query Params: `platform_uid`, `platform_token`
   - Response: 高解像度サムネイルURL、説明文、放送日など

5. **動画ストリーム情報取得**
   - URL: `https://statics.tver.jp/content/episode/{episode_id}.json?v={version}`
   - Response: Streaks/M3U8プレイリストURL

### 認証とDRM対策

- トークンの有効期限管理
- リージョンチェック対策（Forwardedヘッダー）
- Widevine DRM非対応コンテンツのみダウンロード可能

## 5. Go移植アーキテクチャ設計

### パッケージ構成

```bash
tverec-go/
├── cmd/
│   └── tverec/          # CLIエントリーポイント
├── internal/
│   ├── config/          # 設定ファイル読み書き
│   ├── downloader/      # 動画ダウンロード機能
│   ├── parser/          # URL/キーワード解析
│   ├── tver/             # TVer APIクライアント
│   └── utils/           # 共通ユーティリティ
├── pkg/
│   └── ...              # 外部パッケージ向け公開API
├── assets/              # デフォルト設定ファイルなど
├── docs/                # ドキュメント
├── go.mod
└── go.sum
```

### 各パッケージの責務

- **`cmd/tverec`**: CLIインターフェース、ユーザー入力処理、機能呼び出し
- **`internal/config`**: JSON/YAML設定ファイルの読み書き、デフォルト値管理
- **`internal/downloader`**: M3U8プレイリスト解析、並列ダウンロード、プログレス表示
- **`internal/parser`**: URLからコンテンツID抽出、キーワード検索、プレイリスト生成
- **`internal/tver`**: TVer APIクライアント、認証トークン管理、メタデータ取得
- **`internal/utils`**: ログ、ファイル操作、日付処理などの共通機能

## 6. Go移植作業タスクリスト

### Phase 1: 基盤構築

1. プロジェクト構造のセットアップ
   - ディレクトリ構成と基本的なgo.mod作成
   - 各パッケージのディレクトリと空の.goファイル作成

2. 設定管理機能の実装
   - configパッケージ: 設定ファイル読み書き
   - CLIフラグと環境変数の統合

3. ログとユーティリティの実装
   - utilsパッケージ: ログ出力、ファイル操作、日付処理

### Phase 2: APIクライアント実装

1. TVer APIクライアントの実装
   - tverパッケージ: 認証トークン取得機能
   - エピソード/シリーズ情報取得APIの実装
   - APIレスポンスのデータ構造定義

### Phase 3: コア機能実装

1. 動画解析機能の実装
   - parserパッケージ: URLからコンテンツIDを抽出
   - キーワード検索機能の実装
   - プレイリスト情報の解析と整理

2. ダウンローダー機能の実装
   - downloaderパッケージ: 動画フォーマット選択ロジック
   - 並列ダウンロード機能とプログレス表示
   - メタデータとサムネイルの保存機能

### Phase 4: CLIアプリケーション実装

1. メインCLIアプリケーションの実装
   - mainパッケージ: コマンドラインインターフェース
   - ユーザー入力の処理と機能統合
   - ヘルプ表示とエラーハンドリング

### Phase 5: 統合・テスト

1. 統合テストとバグ修正
   - 全体機能の結合テスト
   - エラーケースのハンドリング改善
   - パフォーマンスとリソース使用量の最適化

2. ドキュメント整備
   - README.mdの作成
   - 利用手順書の作成
   - 開発者向けドキュメントの作成

## 7. 将来的な拡張機能

- GUIアプリケーション対応（Fyneまたは他のGo GUIフレームワーク）
- ダウンロード履歴管理と再開機能
- 複数同時ダウンロードキュー
- プラグイン機構による機能拡張