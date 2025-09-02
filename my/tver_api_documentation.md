# TVer API Documentation (Unofficial)

This document outlines the TVer API endpoints as used in the Go prototype.

## 1. Get Platform Token

- **Purpose:** Obtains a platform UID and token required for subsequent API calls.
- **Endpoint:** `POST https://platform-api.tver.jp/v2/api/platform_users/browser/create`
- **Method:** `POST`
- **Headers:**
  - `Content-Type`: `application/x-www-form-urlencoded`
  - `User-Agent`: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`
- **Request Body:** `device_type=pc`
- **Successful Response (JSON):**
  ```json
  {
    "Result": {
      "platform_uid": "...",
      "platform_token": "..."
    }
  }
  ```

## 2. Get Series Seasons

- **Purpose:** Fetches a list of season IDs associated with a specific series.
- **Endpoint:** `GET https://platform-api.tver.jp/service/api/v1/callSeriesSeasons/{seriesID}`
- **Method:** `GET`
- **URL Parameters:**
  - `{seriesID}`: The ID of the series (e.g., `srrazrs5j2`).
- **Query Parameters:**
  - `platform_uid`: The UID obtained from the token endpoint.
  - `platform_token`: The token obtained from the token endpoint.
- **Headers:**
  - `User-Agent`: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`
- **Successful Response (JSON):**
  ```json
  {
    "Result": {
      "Contents": [
        {
          "Type": "season",
          "Content": {
            "Id": "..." // Season ID
          }
        }
      ]
    }
  }
  ```

## 3. Get Season Episodes

- **Purpose:** Fetches a list of episodes for a specific season.
- **Endpoint:** `GET https://platform-api.tver.jp/service/api/v1/callSeasonEpisodes/{seasonID}`
- **Method:** `GET`
- **URL Parameters:**
  - `{seasonID}`: The ID of the season.
- **Query Parameters:**
  - `platform_uid`: The UID obtained from the token endpoint.
  - `platform_token`: The token obtained from the token endpoint.
- **Headers:**
  - `User-Agent`: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`
- **Successful Response (JSON):**
  ```json
  {
    "Result": {
      "Contents": [
        {
          "Type": "episode",
          "Content": {
            "Id": "...", // Episode ID
            "Title": "...",
            "EndAt": 1234567890
          }
        }
      ]
    }
  }
  ```
