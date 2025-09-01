# 設計文書

## 概要

Next.jsベースの太陽系シミュレーションWebアプリケーションです。Three.jsを使用した3D描画、React Contextによる状態管理、Vitestによるテスト環境を構築します。レスポンシブデザインでデスクトップとタブレットの両方に対応します。

## アーキテクチャ

### 技術スタック
- **フレームワーク**: Next.js 14 (App Router)
- **3Dライブラリ**: Three.js + React Three Fiber
- **状態管理**: React Context + useReducer
- **スタイリング**: Tailwind CSS
- **テスト**: Vitest + React Testing Library
- **型安全性**: TypeScript

### アプリケーション構造

src/
├── app/                    # Next.js App Router
│   ├── layout.tsx         # ルートレイアウト
│   └── page.tsx           # メインページ
├── components/            # Reactコンポーネント
│   ├── ui/               # UI基本コンポーネント
│   ├── solar-system/     # 太陽系関連コンポーネント
│   └── controls/         # 制御UI
├── hooks/                # カスタムフック
├── lib/                  # ユーティリティ
├── types/                # TypeScript型定義
└── data/                 # 惑星データ


## コンポーネントとインターフェース

### 主要コンポーネント

#### 1. SolarSystemSimulator
- メインコンテナコンポーネント
- 3Dシーンとコントロールパネルを統合
- シミュレーション状態を管理

#### 2. Scene3D
- Three.js/React Three Fiberによる3D描画
- カメラ制御、ライティング、レンダリング
- 惑星とその軌道の描画

#### 3. Planet
- 個別の惑星コンポーネント
- クリックイベント処理
- テクスチャとマテリアル適用

#### 4. ControlPanel
- 再生/一時停止、速度調整、リセット
- シミュレーション制御UI

#### 5. PlanetInfoPanel
- 選択された惑星の詳細情報表示
- モーダル形式での情報表示

### インターフェース定義

typescript
interface Planet {
  id: string;
  name: string;
  nameJa: string;
  radius: number;
  distance: number;
  orbitalPeriod: number;
  rotationPeriod: number;
  mass: number;
  color: string;
  textureUrl: string;
  description: string;
}

interface SimulationState {
  isPlaying: boolean;
  speed: number;
  time: number;
  selectedPlanet: Planet | null;
}

interface CameraControls {
  position: [number, number, number];
  target: [number, number, number];
  zoom: number;
}


## データモデル

### 惑星データ

実際の天体データに基づいた惑星情報：

- 軌道半径（天文単位）
- 公転周期（地球日）
- 自転周期（時間）
- 質量（地球質量比）
- 物理的特性（半径、色、テクスチャ）

### 状態管理

React ContextとuseReducerを使用：

- SimulationContext: シミュレーション状態
- CameraContext: 3Dカメラ制御
- UIContext: UI状態（パネル表示など）

## エラーハンドリング

### 3D描画エラー

- WebGL非対応ブラウザの検出
- フォールバック表示の提供
- エラーバウンダリによる例外処理

### リソース読み込みエラー

- テクスチャ読み込み失敗時のデフォルト表示
- 段階的な機能低下（Progressive Enhancement）

### パフォーマンス最適化

- LOD（Level of Detail）による描画最適化
- フレームレート監視とパフォーマンス調整
- メモリ使用量の監視

## テスト戦略

### 単体テスト（Vitest）

- ユーティリティ関数のテスト
- 惑星データ計算ロジックのテスト
- 状態管理ロジックのテスト

### コンポーネントテスト（React Testing Library）

- UI コンポーネントの描画テスト
- ユーザーインタラクションのテスト
- 状態変更の検証

### 統合テスト

- シミュレーション全体の動作テスト
- 3D描画とUI連携のテスト

### E2Eテスト考慮事項

- 3D描画の検証は困難なため、主要なユーザーフローに焦点
- 惑星クリック → 情報表示の流れ
- コントロールパネルの操作

## 実装詳細

### 3D描画最適化

- インスタンス化による描画パフォーマンス向上
- シェーダーによるカスタム効果
- アニメーションループの最適化

### レスポンシブ対応

- Tailwind CSSによるレスポンシブレイアウト
- タッチデバイス向けのジェスチャー対応
- 画面サイズに応じたUI調整

### アクセシビリティ

- キーボードナビゲーション対応
- スクリーンリーダー対応
- 色覚異常への配慮

## パフォーマンス考慮事項

### 描画最適化

- フレームレート60fps維持
- GPU使用率の監視
- バッテリー消費の最小化

### メモリ管理

- テクスチャの適切な解放
- 不要なオブジェクトのクリーンアップ
- メモリリークの防止

### 読み込み最適化

- 段階的なリソース読み込み
- プリローディング戦略
- 圧縮テクスチャの使用

