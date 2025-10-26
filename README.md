# Task App LayerX

タスク管理アプリケーション - Clean Architecture実装

## 📁 プロジェクト構成

```
task_app_layerx/
├── backend/           # バックエンド (Go + Echo)
│   ├── cmd/          # アプリケーションエントリーポイント
│   ├── internal/     # 内部パッケージ
│   │   ├── domain/           # ドメイン層
│   │   ├── usecase/          # ユースケース層
│   │   ├── infrastructure/   # インフラ層
│   │   └── presentation/     # プレゼンテーション層
│   ├── pkg/          # 外部公開可能なパッケージ
│   ├── migrations/   # データベースマイグレーション
│   ├── go.mod
│   ├── Makefile
│   ├── Dockerfile
│   └── .env
├── frontend/         # フロントエンド (React + TypeScript + Vite)
│   ├── src/
│   ├── package.json
│   └── Dockerfile
└── docker-compose.yml
```

## 🚀 クイックスタート

### 前提条件

- Docker & Docker Compose
- (ローカル開発の場合) Go 1.23+, Node.js 20+

### Docker で起動（推奨）

```bash
# すべてのサービスを起動
docker compose up -d

# データベースマイグレーション実行（初回のみ必須）
cd backend
make migrate-up
cd ..

# ログを確認
docker compose logs -f

# 停止
docker compose down
```

起動後、以下のURLでアクセス可能:
- **フロントエンド**: http://localhost:5173
- **バックエンドAPI**: http://localhost:8080
- **Adminer (DB管理)**: http://localhost:8081
- **MySQL**: localhost:3306

### ローカル開発

#### バックエンド

```bash
cd backend

# 依存関係のインストール
go mod download

# データベース起動（Dockerのみ）
docker compose up -d db

# マイグレーション実行
make migrate-up

# サーバー起動
make run
```

#### フロントエンド

```bash
cd frontend

# 依存関係のインストール
npm install

# 開発サーバー起動
npm run dev
```

## 🗄️ データベース

### マイグレーション

```bash
cd backend

# マイグレーション実行
make migrate-up

# ロールバック
make migrate-down

# 状態確認
make migrate-status
```

### 接続情報

- **Host**: localhost
- **Port**: 3306
- **Database**: task_db
- **User**: task_user
- **Password**: task_password

## 🔧 開発コマンド

### バックエンド

```bash
cd backend

# アプリケーション起動
make run

# Lint実行
make lint

# テスト実行
make test

# ビルド
go build -o bin/server ./cmd/api
```

### フロントエンド

```bash
cd frontend

# 開発サーバー起動
npm run dev

# ビルド
npm run build

# プレビュー
npm run preview
```

## 🏗️ アーキテクチャ

### Clean Architecture

```
Presentation Layer (HTTP Handler)
        ↓
UseCase Layer (Business Logic)
        ↓
Domain Layer (Entities & Interfaces)
        ↑
Infrastructure Layer (DB, External Services)
```

### 主要な技術スタック

**バックエンド:**
- Go 1.23
- Echo v4 (Web Framework)
- MySQL 8.0
- JWT認証
- Clean Architecture

**フロントエンド:**
- React 19
- TypeScript
- Vite 7
- CSS Modules

## 📝 環境変数

### backend/.env

`.env.example`をコピーして`.env`を作成してください：

```bash
cp backend/.env.example backend/.env
```

**JWT_SECRETの生成:**

以下のコマンドで安全なシークレットキーを生成し、`.env`の`JWT_SECRET`に貼り付けてください：

```bash
openssl rand -base64 32
```


## 🔐 認証

JWT（JSON Web Token）を使用した認証を実装しています。

### ユーザー登録

```bash
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "name": "Test User"
  }'
```

### ログイン

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

レスポンスで返される`token`を以降のAPIリクエストのAuthorizationヘッダーに含めてください:

```bash
Authorization: Bearer <token>
```

## 📚 API エンドポイント

### 認証

- `POST /api/v1/auth/signup` - ユーザー登録
- `POST /api/v1/auth/login` - ログイン
- `POST /api/v1/auth/logout` - ログアウト（要認証）
- `GET /api/v1/users` - ユーザー一覧取得（要認証）

### タスク

- `POST /api/v1/tasks` - タスク作成（要認証）
- `GET /api/v1/tasks` - タスク一覧取得（要認証）
- `GET /api/v1/tasks/:id` - タスク詳細取得（要認証）
- `PATCH /api/v1/tasks/:id` - タスク更新（要認証）
- `DELETE /api/v1/tasks/:id` - タスク削除（要認証）

## 🧪 テスト

```bash
cd backend
make test
```