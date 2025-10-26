# 要件書：簡易タスク管理アプリケーション（Go + Echo + MySQL）

## 結論

クリーンアーキテクチャ × OpenAPI-First、検索（status/due_date/priority のフィルタ＆ソート）付き。CI は GitHub Actions で lint / test / commit メッセージ / OpenAPI 検証まで回す——この前提で正式な要件書を一枚に集約しました。

---

## 1. 目的

- 認証済みユーザーが自分のタスクを作成・共有・管理できる Web API を提供する
- 共同作業のため、タスクのオーナーが複数ユーザーをアサイン可能
- 認可で他人のタスク操作を防止。API 契約は OpenAPI で一元管理

---

## 2. スコープ

### 機能（必須）

1. **ユーザー認証**：サインアップ／ログイン／ログアウト（JWT）
2. **タスク CRUD**：一覧・詳細・作成・更新・削除
3. **アサイン**：タスクに複数ユーザーを紐付け
4. **認可**：
   - オーナー＝参照・更新・削除・アサイン管理
   - アサイン先＝参照のみ（編集不可）
   - 無関係者＝不可
5. **エラーハンドリング**：統一 JSON 形式
6. **ドキュメント**：OpenAPI（Swagger UI）＋ README（起動・使用方法）

### 非機能（必須）

- Clean Architecture 準拠（Presentation / UseCase / Domain / Infrastructure）
- OpenAPI-First（oapi-codegen で型＆Echo サーバ IF 生成）
- Docker Compose（api / db / adminer）
- Migration（golang-migrate）
- CI（GitHub Actions）：lint / vet / build / test / commitlint / OpenAPI validate
- ログ／リカバリ／CORS／Request-ID ミドルウェア

### スコープ外（今回含めない）

- パスワードリセット／メール検証
- 通知・Webhook・フロントエンド UI
- デプロイ（CD は無し）
- **検索機能**（一覧のフィルタ・ソート・ページング）※今後実装予定

---

## 3. 使用技術・ツール

- Go 1.22+, Echo v4 / MySQL 8.x
- Auth: JWT(HS256), bcrypt
- OpenAPI v3 + oapi-codegen
- Lint/Test: golangci-lint, go test(+race), testify, sqlmock, httptest
- docker / docker-compose, golang-migrate, Adminer

---

## 4. アーキテクチャ

```
Presentation (Echo, oapi validator, DTO変換)
      ↓
UseCase (アプリケーションサービス/認可/Tx境界)
      ↓
Domain (Entity/Value/Policy)
      ↑
Infrastructure (MySQL Repository, TxManager, Clock)
```

### 代表ディレクトリ

```
/api/openapi.yaml
/cmd/api/main.go
/internal/presentation/http/echo/   # oapi-server 実装, middleware
/internal/presentation/converter/   # OAPI <-> Domain
/internal/usecase/                  # ports & usecases
/internal/domain/                   # entities, policies
/internal/infrastructure/mysql/     # repositories (sql)
/internal/infrastructure/tx/        # TxManager
/pkg/auth/jwt, /pkg/hash/bcrypt, /pkg/config
/migrations/*.sql
/.github/workflows/ci.yml
```

---

## 5. データモデル（ER）

```mermaid
erDiagram
    users {
      bigint id PK
      varchar email UNIQUE NOT NULL
      varchar password_hash NOT NULL
      varchar name
      int token_version DEFAULT 0
      datetime created_at
      datetime updated_at
      datetime deleted_at
    }

    tasks {
      bigint id PK
      bigint owner_id FK
      varchar title NOT NULL
      text description
      datetime due_date
      enum status
      tinyint priority DEFAULT 0
      datetime created_at
      datetime updated_at
      datetime deleted_at
      INDEX idx_owner
      INDEX idx_due
      INDEX idx_status
      INDEX idx_priority
    }

    task_assignees {
      bigint task_id FK
      bigint user_id FK
      bigint assigned_by FK
      datetime created_at
      PK task_id_user_id
      INDEX idx_user
    }

    users ||--o{ tasks : owns
    tasks ||--o{ task_assignees : has
    users ||--o{ task_assignees : assigns
```

### テーブル詳細

#### users
- `id`: BIGINT, PK
- `email`: VARCHAR(255), UNIQUE, NOT NULL
- `password_hash`: VARCHAR(255), NOT NULL
- `name`: VARCHAR(100)
- `token_version`: INT, DEFAULT 0
- `created_at`, `updated_at`, `deleted_at`: DATETIME

#### tasks
- `id`: BIGINT, PK
- `owner_id`: BIGINT, FK → users.id
- `title`: VARCHAR(255), NOT NULL
- `description`: TEXT
- `due_date`: DATETIME
- `status`: ENUM('TODO', 'IN_PROGRESS', 'DONE'), DEFAULT 'TODO'
- `priority`: TINYINT, DEFAULT 0
- `created_at`, `updated_at`, `deleted_at`: DATETIME
- INDEX: `idx_owner(owner_id)`, `idx_due(due_date)`, `idx_status(status)`, `idx_priority(priority)`

#### task_assignees
- `task_id`: BIGINT, FK → tasks.id
- `user_id`: BIGINT, FK → users.id
- `assigned_by`: BIGINT, FK → users.id
- `created_at`: DATETIME
- PK: (task_id, user_id)
- INDEX: `idx_user(user_id)`

---

## 6. 認証・認可

### 認証

- **サインアップ**：email unique、password は bcrypt（cost=12 目安）
- **ログイン**：一致で JWT 発行（sub=userID, tv=token_version, exp=1h）
- **ログアウト**：token_version++ で既存トークン一括失効

### 認可ポリシ

- **参照**：owner == me または assignee == me
- **更新/削除/アサイン変更**：owner == me のみ

---

## 7. API 仕様（OpenAPI v3 要旨）

- **ルート**：`/api/v1`
- **セキュリティ**：bearerAuth（JWT）

### 主要エンドポイント

| メソッド | パス | 機能 | 認証 |
|---------|------|------|------|
| POST | `/auth/signup` | ユーザー登録 | 公開 |
| POST | `/auth/login` | ログイン | 公開 |
| POST | `/auth/logout` | ログアウト | 要JWT |
| GET | `/tasks` | タスク一覧（検索込み） | 要JWT |
| POST | `/tasks` | タスク作成 | 要JWT |
| GET | `/tasks/{id}` | タスク詳細 | 要JWT |
| PATCH | `/tasks/{id}` | タスク更新 | 要JWT（ownerのみ） |
| DELETE | `/tasks/{id}` | タスク削除 | 要JWT（ownerのみ） |

### /tasks の検索パラメータ

- `status`：`TODO` | `IN_PROGRESS` | `DONE`
- `dueFrom`, `dueTo`：ISO8601 datetime
- `priorityMin`, `priorityMax`：0..5
- `sort`：カンマ区切り。許可列＝`dueDate` | `priority` | `createdAt`、`-` で降順
  - 例：`sort=-dueDate,priority`
- `page`（default 1）, `perPage`（default 20, max 100）

### エラーフォーマット（統一）

```json
{
  "code": "VALIDATION_ERROR",
  "message": "title is required",
  "details": { "field": "title" }
}
```

---

## 8. 検索実装の要点

- **WHERE 句**：`owner_id = me OR assignees.user_id = me` に限定（未関係非表示）
- **フィルタは任意適用**（NULL なら無視）
- **ソートはホワイトリスト正規化**：
  - マップ：`{"dueDate":"due_date", "priority":"priority", "createdAt":"created_at"}`
  - `-` 有無で ASC/DESC 決定
  - 未指定は `created_at DESC, id DESC`
- **ページング**：`LIMIT ? OFFSET ?`
- **カウント**：`COUNT(*)` 別クエリ or `FOUND_ROWS()`

---

## 9. マイグレーション（要点）

- `users`(email UNIQUE, password_hash, token_version, timestamps, deleted_at)
- `tasks`(owner_id FK, title, description, due_date, status, priority, indexes, soft-delete)
- `task_assignees`(task_id, user_id, assigned_by, PK(task_id, user_id))
- 追加インデックス：status, priority, 複合 (owner_id, status, due_date)（任意）

---

## 10. ミドルウェア

- Recovery / Logger / RequestID / CORS
- JWT 認証（Authorization: Bearer）
- oapi RequestValidator（スキーマ準拠チェック）
- 例外変換：ドメインエラー→HTTP ステータス＋標準エラー JSON

---

## 11. CI（GitHub Actions）

### トリガ

- PR（main, develop）／main への push

### ジョブ

#### 1. build-lint-test
- `go build ./...` / `go vet ./...` / `golangci-lint` / `go test ./... -race -cover`
- Go モジュール＆ビルドキャッシュ

#### 2. commitlint（PR 時）
- Conventional Commits ルール検査（`feat:`, `fix:`, ...）

#### 3. openapi-validate
- `api/openapi.yaml` バリデーション（スキーマ整合性）

### 設定ファイル

- `.github/workflows/ci.yml`
- `.golangci.yml`
- `.commitlintrc.json`

---

## 12. テスト

- **Unit**：ドメイン（認可ポリシ、値検証）
- **UseCase**：リポジトリモックで分岐（owner/assignee/others、403/404/200）
- **Repository**：sqlmock（成功/異常）、必要に応じ dockertest
- **Handler**：httptest（JWT 有無、Validator 経由の 400、403/404/200）
- **カバレッジ**：CI で収集（artifact）

---

## 13. 運用・起動

### Docker Compose

- `api:8080` / `mysql:3306` / `adminer:8081`

### .env 例

```env
APP_PORT=8080
DB_HOST=db
DB_PORT=3306
DB_USER=app
DB_PASS=secret
DB_NAME=taskapp
JWT_SECRET=change-me
JWT_EXPIRES=3600
```

### 代表コマンド

- `make up`（compose up -d）
- `make migrate-up`
- `make gen-oapi`
- `make run`
- Swagger: `GET /swagger/index.html`

---

## 14. 受け入れ基準（Acceptance Criteria）

### 認証
- サインアップ→ログイン→JWT で認証済みエンドポイントにアクセスできる
- ログアウト後、古い JWT は 401

### タスク
- オーナーは CRUD & アサイン可能
- アサイン先は参照でき、編集しようとすると 403
- 無関係者は 404（存在隠蔽）

### 検索
- 各フィルタが正しく作用
- ソート：`-dueDate,priority` の順序で返る
- ページング：total/page/perPage/items の整合

### その他
- エラーは統一 JSON 形式
- CI：PR 作成で lint/test/commitlint/openapi がすべてパス

---

## 15. 納品物

- ソースコード（リポジトリ）
- `api/openapi.yaml`
- `migrations/*.sql`
- `.github/workflows/ci.yml`, `.golangci.yml`, `.commitlintrc.json`
- `README.md`（セットアップ、起動、API 利用例、検索クエリ例、CI 説明）
- （任意）テストレポート／カバレッジファイル

---

## 16. 例：検索 API 利用

```bash
GET /api/v1/tasks?status=IN_PROGRESS&dueFrom=2025-10-01T00:00:00Z&dueTo=2025-10-31T23:59:59Z&priorityMin=2&sort=-dueDate,priority&page=1&perPage=20
Authorization: Bearer <token>
```

---

## 17. 次アクション（実装順）

1. `api/openapi.yaml` をこの要件で確定
2. oapi-codegen 生成 & Echo 起動（JWT / Validator 連携）
3. UseCase / Repository 骨組み & ListForUser（検索）実装
4. マイグレーション適用・インデックス確認
5. CI ファイル配置・初回 PR でパス確認
6. README 仕上げ

---

必要なら、この要件書に対応する OpenAPI 初版 YAML と CI YAML の実体も続けて貼れます。どこから着手するかだけ指示ください（その章を即ペーストします）。