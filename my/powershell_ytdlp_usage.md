# PowerShell版 TVerRecにおけるyt-dlp利用方法

このドキュメントは、TVerRecのPowerShell版が`yt-dlp`（またはそのフォーク）をどのように利用しているかをまとめたものです。

## 1. 中心的な役割

`yt-dlp`は、TVerからの動画ダウンロード、メタデータ・字幕・サムネイルの取得と埋め込みという、プロジェクトの中核的な処理を担っています。

## 2. 主な実行関数

`yt-dlp`の呼び出しは、主に`src/functions/tverrec_functions.ps1`に定義されている以下の2つのラッパー関数によって行われます。

-   `Invoke-Ytdl`: TVerの動画をダウンロードするためのメイン関数。
-   `Invoke-NonTverYtdl`: TVer以外のサイトからダウンロードするための関数。

## 3. `Invoke-Ytdl`で利用される主なコマンドライン引数

以下は、`Invoke-Ytdl`関数内で動的に構築される`yt-dlp`の主な引数です。

-   **出力とフォーマット:**
    -   `--output "{episodeID}.{ext}"`: 出力ファイル名を一時的にエピソードID（例: `epxxxxxxxx.mp4`）に設定します。
    -   `--merge-output-format mp4`: ダウンロードしたストリームをMP4コンテナにマージします。
    -   `--paths`: 一時ファイルと最終的な動画ファイルの保存先ディレクトリを指定します。

-   **メタデータと字幕:**
    -   `--embed-thumbnail`: 番組のサムネイル画像を動画ファイルに埋め込みます。
    -   `--embed-chapters`: チャプター情報を埋め込みます。
    -   `--embed-metadata`: 番組名、シリーズ名などのメタデータを埋め込みます。
    -   `--embed-subs`: 字幕データを埋め込みます。
    -   `--sub-langs all`: 利用可能なすべての言語の字幕を取得します。
    -   `--convert-subs srt`: 字幕をSRT形式に変換します。

-   **ネットワークとパフォーマンス:**
    -   `--concurrent-fragments <N>`: 動画のフラグメント（断片）を並列でダウンロードし、高速化を図ります。
    -   `--limit-rate <RATE>`: ダウンロード速度を制限します。
    -   `--geo-verification-proxy <URL>`: Geo-IPチェックを回避するためにプロキシサーバーを利用します。
    -   `--xff <IP>`: `X-Forwarded-For`ヘッダーを送信し、IPアドレスを偽装するオプションです。

-   **外部ツール連携:**
    -   `--ffmpeg-location <PATH>`: `ffmpeg`の実行ファイルパスを指定します。
    -   `--exec "after_video:{command}"`: ダウンロード完了後に指定したPowerShellコマンドを実行します。主に、一時的なファイル名（`epxxxxxxxx.mp4`）から最終的なフォーマットのファイル名への変更（リネーム）に利用されます。

## 4. プロセス管理

-   **タイムアウト処理:** `Start-Job`を利用して`yt-dlp`プロセスを監視し、設定されたタイムアウト時間 (`$script:ytdlTimeoutSec`) を超えた場合にプロセスを強制終了させる仕組みが実装されています。
-   **並列ダウンロード制御:** `Get-YtdlProcessCount`関数で実行中の`yt-dlp`プロセス数を監視し、設定された最大並列数を上回らないようにダウンロードの開始を待機する制御が行われています。
