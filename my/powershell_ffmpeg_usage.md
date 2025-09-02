# PowerShell版 TVerRecにおけるffmpeg利用方法

このドキュメントは、TVerRecのPowerShell版が`ffmpeg`および`ffprobe`をどのように利用しているかをまとめたものです。

## 1. 主な役割

`ffmpeg`は、主に2つの役割を担っています。

1.  **動画ファイルの整合性チェック:** ダウンロードした動画ファイルが破損していないか検証するために使用されます。
2.  **`yt-dlp`の依存ツール:** `yt-dlp`が動画と音声ストリームを結合（mux）したり、フォーマットを変換したりする際に、内部的に`ffmpeg`を呼び出します。

## 2. 動画の整合性チェック (`Invoke-IntegrityCheck`)

動画の検証は`src/functions/tverrec_functions.ps1`内の`Invoke-IntegrityCheck`関数で行われ、設定に応じて2つのモードがあります。

### a. 完全検査モード (`ffmpeg`)

`simplifiedValidation`が`$false`の場合、`ffmpeg`を使用して動画全体をデコードし、エラーがないか徹底的にチェックします。

-   **実行コマンド (例):**
    ```powershell
    ffmpeg -hide_banner -v error -xerror -i "path/to/video.mp4" -f null -
    ```
-   **主な引数:**
    -   `-v error`: ログレベルをエラーのみに設定し、不要な出力を抑制します。
    -   `-xerror`: デコード中にエラーが発生した場合、即座にプロセスを終了させます。
    -   `-i "{video_path}"`: 入力ファイルとして検証対象の動画を指定します。
    -   `-f null -`: 出力フォーマットを`null`に設定し、ファイルへの書き込みを行わずにデコード処理のみを実行します。これにより高速な検証が可能です。
    -   **ハードウェアアクセラレーション:** ユーザーは設定 (`$script:ffmpegDecodeOption`) で`-hwaccel qsv`などのデコードオプションを追加できます。

### b. 簡易検査モード (`ffprobe`)

`simplifiedValidation`が`$true`の場合、`ffprobe`を使用してより高速な簡易チェックを行います。

-   **実行コマンド (例):**
    ```powershell
    ffprobe -hide_banner -v error -err_detect explode -i "path/to/video.mp4"
    ```
-   **主な引数:**
    -   `-err_detect explode`: ストリームの構文エラーなどを検知した場合、即座に処理を中断します。

## 3. プロセス実行とエラー判定

-   `Invoke-FFmpegProcess`関数が`ffmpeg`/`ffprobe`のプロセスを起動します。
-   標準エラー出力がログファイル (`ffmpeg_err_*.log`) にリダイレクトされます。
-   検証の成否は、プロセスの**終了コード**と、エラーログの**行数**の両方を見て判断されます。エラー行が多数（30行以上）ある場合も失敗と見なされます。

## 4. `yt-dlp`からの利用

`yt-dlp`は、動画と音声を結合したり、字幕を埋め込んだりするために`ffmpeg`を必要とします。TVerRecは`yt-dlp`を呼び出す際に、`--ffmpeg-location`引数を使って`ffmpeg`の実行ファイルのパスを明示的に伝えています。
