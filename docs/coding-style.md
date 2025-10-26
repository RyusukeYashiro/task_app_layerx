# コーディング規約（Task Management API）

本ドキュメントは、Task Management API プロジェクトの開発における最小限のコーディングルールを定義する。

すべての開発者はこの規約を遵守し、チーム全体で一貫したコード品質・可読性・保守性を維持することを目的とする。

---

## 1. 命名規則

### 基本方針
- **Initialism は大文字**: `ID`, `URL`, `HTTP`, `JWT`, `DB` など
  - ✅ `UserID`, `FindByID`
  - ❌ `UserId`, `FindById`
- **可視性**: 外部に出さない要素は先頭小文字（パッケージ内専用）
- **関数は動詞始まり**: `CreateTask`, `UpdateUser`, `FindByID`
- **構造体は単数名詞**: `User`, `Task`, `Project`
- **JSONキーは camelCase に統一**: `createdAt`, `dueDate`, `assigneeIds`
- **空スライスは `[]` を返す**: `null` や省略は不可

---

## 2. GORM モデル設計

### 基本方針
- **GORM v2 前提**
- **ソフトデリート**: `gorm.DeletedAt` フィールドを必ず持つ
  - `CreatedAt`, `UpdatedAt`, `DeletedAt` を標準装備
- **NULL 許容**: ポインタ型を使用（`*string`, `*time.Time` など）
  - OpenAPIの `nullable` に対応
- **外部キー/インデックス**: タグで明示（実DBマイグレーションと整合）
  - 例: `gorm:"foreignKey:OwnerID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
  - 例: `gorm:"index:idx_owner"`
- **多対多**: `many2many:<join_table>` を使用
  - 例: `gorm:"many2many:task_assignees"`
- **テーブル名**: 必要時のみ `TableName()` メソッドで明示

### モデル例

```go
type User struct {
    ID           int64          `gorm:"primaryKey"`
    Email        string         `gorm:"uniqueIndex;not null"`
    PasswordHash string         `gorm:"not null"`
    Name         string         `gorm:"not null"`
    TokenVersion int            `gorm:"default:0"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
    DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type Task struct {
    ID          int64          `gorm:"primaryKey"`
    OwnerID     int64          `gorm:"not null;index:idx_owner"`
    Title       string         `gorm:"not null"`
    Description *string
    DueDate     *time.Time     `gorm:"index:idx_due"`
    Status      string         `gorm:"type:enum('TODO','IN_PROGRESS','DONE');default:'TODO';index:idx_status"`
    Priority    int8           `gorm:"default:0;index:idx_priority"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   gorm.DeletedAt `gorm:"index"`
    
    Owner     User   `gorm:"foreignKey:OwnerID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
    Assignees []User `gorm:"many2many:task_assignees"`
}
```

---

## 3. トランザクション & DB 操作

### 基本方針
- **書き込み系は必ず Tx 内**
  - 基本形: `db.WithContext(ctx).Transaction(func(tx *gorm.DB) error { ... })`
- **Context 必須**: DBアクセスは常に `WithContext(ctx)`
- **スコープ関数で条件を合成**: 認可/検索条件/期間/優先度など
- **ソートはホワイトリスト**: 許可された列名のみ
- **ループ内クエリ禁止**: 必要なら一括取得→メモリ結合

### トランザクション例

```go
// UseCase層でトランザクション境界を定義
func (u *TaskUseCase) CreateTask(ctx context.Context, req CreateTaskRequest) (*Task, error) {
    var task *domain.Task
    
    err := u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // 1. タスク作成
        task = &domain.Task{...}
        if err := tx.Create(task).Error; err != nil {
            return err
        }
        
        // 2. アサイン情報追加
        if len(req.AssigneeIDs) > 0 {
            if err := tx.Model(task).Association("Assignees").Append(assignees); err != nil {
                return err
            }
        }
        
        return nil
    })
    
    return task, err
}
```

### スコープ関数例

```go
// 認可スコープ: 自分がオーナーまたはアサインされているタスクのみ
func ScopeAccessibleTasks(userID int64) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("owner_id = ?", userID).
            Or("id IN (SELECT task_id FROM task_assignees WHERE user_id = ?)", userID)
    }
}

// 使用例
db.WithContext(ctx).
    Scopes(ScopeAccessibleTasks(currentUserID)).
    Find(&tasks)
```

---

## 4. エラーハンドリング & ログ

### 基本方針
- **error 文字列**: 先頭小文字・句点なし
  - 例: `fmt.Errorf("invalid token")`
  - 例: `fmt.Errorf("user not found: %w", err)`
- **エラーラップ**: `%w` で wrap、`errors.Is` / `errors.As` で判定
- **GORM エラー**: `gorm.ErrRecordNotFound` は 404 にマッピング
- **握りつぶし禁止**: ただしユーザー影響のない副作用のみ可（コメントで意図を明記）
- **ログは UseCase層で**: Infrastructure層では返却を優先

### エラーハンドリング例

```go
// ❌ 悪い例
func (r *TaskRepository) FindByID(ctx context.Context, id int64) (*Task, error) {
    var task Task
    if err := r.db.First(&task, id).Error; err != nil {
        return nil, err  // そのまま返すのは NG
    }
    return &task, nil
}

// ✅ 良い例
func (r *TaskRepository) FindByID(ctx context.Context, id int64) (*Task, error) {
    var task Task
    err := r.db.WithContext(ctx).First(&task, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, domain.ErrTaskNotFound  // ドメインエラーに変換
        }
        return nil, fmt.Errorf("failed to find task: %w", err)
    }
    return &task, nil
}
```

### ログ出力

```go
// UseCase層でログ出力
func (u *TaskUseCase) DeleteTask(ctx context.Context, taskID, userID int64) error {
    task, err := u.taskRepo.FindByID(ctx, taskID)
    if err != nil {
        u.logger.Error("failed to find task", "taskID", taskID, "error", err)
        return err
    }
    
    // ...
}
```

---

## 5. JSON レスポンス / API 設計

### 基本方針
- **OpenAPI-First**: `api/openapi.yaml` を唯一の契約とし、生成物と整合
- **成功レスポンス**:
  - `201 Created`: 作成
  - `200 OK`: 更新・取得
  - `204 No Content`: 削除（ボディなし）
- **失敗レスポンス**: `{code, message, details}` 形式で統一（OpenAPI の `ErrorResponse` 準拠）

### 認可規約
- **オーナー**: 参照・更新・削除・アサイン管理
- **アサイン先**: 参照のみ
- **無関係者**: 存在隠蔽（404）

### エラーレスポンス例

```json
{
  "code": "VALIDATION_ERROR",
  "message": "title is required",
  "details": {
    "field": "title"
  }
}
```

### HTTPステータスコードマッピング

| ドメインエラー | HTTPステータス | エラーコード |
|--------------|---------------|-------------|
| `ErrUnauthorized` | 401 | `UNAUTHORIZED` |
| `ErrForbidden` | 403 | `FORBIDDEN` |
| `ErrNotFound` | 404 | `NOT_FOUND` |
| `ErrConflict` | 409 | `CONFLICT` |
| バリデーションエラー | 400 | `VALIDATION_ERROR` |
| その他 | 500 | `INTERNAL_ERROR` |

---

## 6. クリーンアーキテクチャの責務

### 各層の責務

#### Presentation層
- HTTP リクエスト/レスポンス処理
- バリデーション
- DTO変換（OpenAPI型 ↔ Domain型）
- エラー → HTTPステータス変換

#### UseCase層
- アプリケーションフロー
- 認可判定
- トランザクション境界

#### Domain層
- エンティティ
- 値オブジェクト
- ポリシー（CanView, CanEdit等）
- リポジトリインターフェース
- **外部技術非依存**（GORM, Echo, JWT等を直接使用しない）

#### Infrastructure層
- GORM リポジトリ実装
- DB接続
- トランザクション管理

#### pkg層
- 共通技術（JWT, bcrypt, config など）の薄いラッパ

### 依存方向

```
Presentation → UseCase → Domain ← Infrastructure
                            ↑
                           pkg
```

---

## 7. パフォーマンス & セキュリティ

### パフォーマンス
- **N+1 回避**: `Preload` は必要最小限。列は `Select` で絞る
  - ❌ `db.Preload("Owner").Preload("Assignees.User").Find(&tasks)`
  - ✅ 必要な関連のみ Preload
- **ループ内クエリ禁止**: 一括取得→メモリ結合
- **ページネーション必須**: 大量データ取得時

### セキュリティ
- **ユーザー入力はバインド変数**: 生SQLは `?` を必須
  - ❌ `db.Where("status = " + userInput)`
  - ✅ `db.Where("status = ?", userInput)`
- **検索**: 認可条件（owned or assigned）を先に適用
- **JWT**: HS256 固定、`exp`/`iat` 必須、`token_version` で世代照合
- **パスワード**: bcrypt cost=12 目安。平文保存禁止

### N+1回避例

```go
// ❌ 悪い例: N+1問題
func (r *TaskRepository) ListWithOwners(ctx context.Context) ([]*Task, error) {
    var tasks []*Task
    r.db.WithContext(ctx).Find(&tasks)
    
    for i := range tasks {
        r.db.First(&tasks[i].Owner, tasks[i].OwnerID)  // N回クエリ発行
    }
    return tasks, nil
}

// ✅ 良い例: Preload
func (r *TaskRepository) ListWithOwners(ctx context.Context) ([]*Task, error) {
    var tasks []*Task
    err := r.db.WithContext(ctx).
        Preload("Owner").
        Find(&tasks).Error
    return tasks, err
}
```

---

## 8. Lint / CI

### golangci-lint 設定
- **使用 linters**: `govet`, `errcheck`, `staticcheck`, `gosimple`, `unused`, `ineffassign`, `typecheck`, `gofumpt`, `revive`, `misspell`, `dupl`, `errorlint`, `gosec`
- **自動生成ファイル**: `*.gen.go` は lint 対象外
- **import 並び**: `goimports`/`gci`（標準→サードパーティ→自プロジェクト）

### CI フロー
1. `build`: `go build ./...`
2. `lint`: `golangci-lint run`
3. `test`: `go test ./... -race -cover`
4. `openapi-validate`: `api/openapi.yaml` バリデーション
5. `commitlint`: Conventional Commits準拠チェック

### Conventional Commits 形式
```
feat: add task search functionality
fix: resolve N+1 query in task list
docs: update API documentation
test: add unit tests for task usecase
refactor: extract authorization logic to policy
```

---

## 9. 例外ルール（やらないこと）

### 禁止事項
1. **ドメイン層から外部ライブラリ直参照**: JWT/GORM/Echo 禁止
2. **破壊的 API 変更**: 段階移行（追加→フロント移行→旧削除）
3. **null 配列の返却**: 空配列 `[]` を返す
4. **JSONキーの snake_case**: camelCase に統一
5. **パスワード平文保存**: bcrypt必須
6. **エラーの握りつぶし**: 必ずログか返却

---

## 10. 変更手順（API/スキーマ）

### 変更フロー
1. **OpenAPI 更新**: `api/openapi.yaml` を編集
2. **生成コード更新**: `make gen-oapi` 実行
3. **Domain/UseCase/Infrastructure を追従**: 型変更に対応
4. **Handler/Converter 更新**: Presentation層の修正
5. **テスト更新**: 単体テスト・統合テスト追加/修正
6. **CI green 確認**: すべてのチェックが通ることを確認

### 後方互換性の保ち方

```yaml
# ❌ 悪い例: フィールド削除（破壊的変更）
# Before
TaskResponse:
  properties:
    id: ...
    oldField: ...  # これを削除したい

# After
TaskResponse:
  properties:
    id: ...
    # oldField 削除 → クライアント壊れる

# ✅ 良い例: 段階的移行
# Step 1: 新フィールド追加
TaskResponse:
  properties:
    id: ...
    oldField: ...        # まだ残す
    newField: ...        # 追加

# Step 2: クライアント移行期間

# Step 3: 旧フィールド削除
TaskResponse:
  properties:
    id: ...
    newField: ...
```

---

## 11. テスト規約

### テスト種別
- **Unit**: ドメインロジック（認可ポリシ、値検証）
- **UseCase**: リポジトリモックで分岐テスト
- **Repository**: sqlmock または dockertest
- **Handler**: httptest（JWT有無、403/404/200）

### テストファイル配置
```
internal/
  domain/
    policy/
      task_policy.go
      task_policy_test.go  # 同一パッケージ
  usecase/
    task/
      task_usecase.go
      task_usecase_test.go
```

### モックの作成
```go
// UseCase層のテストではRepositoryをモック化
type mockTaskRepository struct {
    mock.Mock
}

func (m *mockTaskRepository) FindByID(ctx context.Context, id int64) (*domain.Task, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.Task), args.Error(1)
}
```

---

## 12. ドキュメント更新

### 必須ドキュメント
- **OpenAPI仕様**: `api/openapi.yaml`（常に最新を保つ）
- **README.md**: セットアップ手順、API使用例
- **CHANGELOG.md**: 機能追加・修正履歴（セマンティックバージョニング）
- **Architecture Decision Records (ADR)**: 重要な設計判断を記録（任意）

### README.md 必須項目
1. プロジェクト概要
2. 技術スタック
3. セットアップ手順
4. 起動方法
5. API エンドポイント一覧
6. 認証方法
7. 開発ガイド
8. テスト実行方法

---

本コーディング規約は、Clean Architecture + GORM v2 + OpenAPI-First 構成での最小限のチーム開発標準です。