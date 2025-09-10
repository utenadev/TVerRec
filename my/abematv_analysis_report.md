# AbemaTV Extractor Analysis and Comparison with TVer

## Introduction

This report analyzes the `yt-dlp` extractor for AbemaTV (`abematv.py`) to determine its functionality and complexity, particularly in comparison to the extractor for TVer (`tver.py`). The goal is to understand if AbemaTV can be targeted for video downloads in a similar fashion to TVer.

## TVer Extractor (`tver.py`) Analysis

The TVer extractor is relatively straightforward. Its key characteristics are:

*   **Reliance on Third-Party Platforms:** TVer does not host its videos directly. Instead, it uses major video platforms like **Brightcove** and **Streaks**. The `yt-dlp` extractor for TVer primarily acts as a bridge.
*   **Simple API Interaction:** The extractor calls TVer's public APIs (`platform-api.tver.jp`, `service-api.tver.jp`) to fetch metadata about episodes and series.
*   **Delegation of DRM/Extraction:** Once it retrieves the video reference ID, it delegates the actual video extraction and handling of DRM to other, more generic InfoExtractors within `yt-dlp` (e.g., `BrightcoveNewIE`, `StreaksBaseIE`). The `tver.py` code itself does not contain any complex video decryption logic.
*   **Simple Authentication:** It establishes a simple, non-user-specific session to communicate with the API. It does not handle user logins.

In essence, targeting TVer is made simple because it builds upon standard, well-supported video distribution platforms.

## AbemaTV Extractor (`abematv.py`) Analysis

The AbemaTV extractor is significantly more complex and demonstrates a much higher level of custom engineering.

*   **Proprietary Infrastructure:** AbemaTV uses its own API (`api.abema.io`), its own license server (`license.abema.io`), and a custom-secured CDN setup. It does not use a standard third-party video platform.
*   **Complex Authentication:** The extractor must first generate a unique Device ID. It then uses a hardcoded secret key (`_SECRETKEY`) and a complex, time-based HMAC hashing algorithm (`_generate_aks`) to create an application key. This key is used to authorize the device and obtain a user token, which is then exchanged for a temporary media token required for playback. It also fully supports user email/password logins for premium content.
*   **Custom DRM Implementation:** This is the most critical difference. AbemaTV uses a **proprietary, custom-built DRM system**. The `yt-dlp` extractor contains a complete, reverse-engineered implementation of this DRM scheme. The process involves:
    1.  A custom protocol handler (`abematv-license://`) to intercept license requests from the video manifest.
    2.  Requesting a license from AbemaTV's license server.
    3.  The server returns an encrypted video key. This key is encoded using a **custom Base58-style alphabet**.
    4.  The extractor decodes this key.
    5.  It then decrypts the key using **AES-ECB**. The decryption key for this step is itself generated via an HMAC-SHA256 hash involving another hardcoded key (`_HKEY`), the device ID, and a `cid` from the license response.
    6.  The final decrypted 16-byte key is then used to download the video stream.

## Conclusion

**Yes, AbemaTV is a targetable platform for `yt-dlp`, but it is fundamentally different from and more complex than TVer.**

*   **TVer:** Easy to support because it uses standard, external video platforms. The extractor is simple and maintainable.
*   **AbemaTV:** Difficult to support because it uses a completely custom and proprietary DRM system. The `yt-dlp` extractor is a sophisticated piece of reverse engineering.

This means that the AbemaTV extractor is **fragile**. If AbemaTV decides to change any part of its custom DRM (e.g., update the secret keys, change the hashing algorithm, alter the license protocol), the extractor will immediately break and will require significant effort to reverse-engineer and fix. TVer is more resilient to such changes as long as its underlying platform providers (Brightcove, Streaks) remain supported by `yt-dlp`.
