# TODO管理アプリケーション

シンプルで使いやすいTODO管理アプリケーションです。

## 技術スタック

### フロントエンド
- **Next.js 14** - Reactフレームワーク
- **TypeScript** - 型安全な開発
- **Apollo Client** - GraphQLクライアント

### バックエンド
- **Go 1.21** - バックエンド言語
- **Gin** - Webアプリケーションフレームワーク
- **JWT** - 認証トークン

### データベース・GraphQL
- **PostgreSQL 15** - リレーショナルデータベース
- **Hasura** - GraphQLエンジン

### インフラ
- **Docker & Docker Compose** - コンテナ化

## 主な機能

### ユーザー機能
- ユーザー登録（メールアドレス・パスワード）
- ログイン/ログアウト
- TODOのCRUD操作
  - 作成
  - 一覧表示
  - 更新（タイトル、説明、完了状態）
  - 削除

### 管理者機能
- ユーザー一覧表示
- ユーザーロール変更（user/admin）
- ユーザー削除
- 全ユーザーのTODO一覧表示

## プロジェクト構造

```
todo-app/
├── backend/               # Goバックエンド
│   ├── cmd/
│   │   └── server/       # メインアプリケーション
│   ├── internal/
│   │   ├── auth/         # 認証ロジック（JWT、パスワード）
│   │   ├── config/       # 設定管理
│   │   ├── handler/      # HTTPハンドラー
│   │   ├── middleware/   # ミドルウェア
│   │   └── model/        # データモデル
│   └── go.mod
├── frontend/              # Next.jsフロントエンド
│   ├── src/
│   │   ├── app/          # Next.js App Router
│   │   ├── components/   # Reactコンポーネント
│   │   ├── lib/          # ユーティリティ関数
│   │   └── types/        # TypeScript型定義
│   └── package.json
├── hasura/               # Hasura設定
│   ├── migrations/       # データベースマイグレーション
│   ├── metadata/         # Hasuraメタデータ
│   └── config.yaml
├── docker/               # Docker設定ファイル
│   ├── backend/
│   │   └── Dockerfile    # バックエンド用Dockerfile
│   ├── frontend/
│   │   └── Dockerfile    # フロントエンド用Dockerfile
│   └── docker-compose.yml # Docker Compose設定
├── .env.example          # 環境変数テンプレート
└── README.md
```

## セットアップ

### 前提条件

- Docker Desktop がインストールされていること
- Docker Compose が利用可能であること

### 1. リポジトリのクローン

```bash
git clone <repository-url>
cd todo-app
```

### 2. 環境変数の設定

```bash
cp .env.example .env
```

必要に応じて `.env` ファイルを編集してください：

```env
# Database
POSTGRES_DB=todo_db
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres

# Hasura
HASURA_GRAPHQL_ADMIN_SECRET=hasura_admin_secret

# JWT Secret (本番環境では必ず変更してください)
JWT_SECRET=your-256-bit-secret-key-change-this-in-production

# Backend
GIN_MODE=debug

# Frontend
NEXT_PUBLIC_HASURA_ENDPOINT=http://localhost:8080/v1/graphql
NEXT_PUBLIC_BACKEND_URL=http://localhost:8000
```

### 3. アプリケーションの起動

```bash
docker compose -f docker/docker-compose.yml up -d
```

初回起動時は、イメージのビルドに数分かかる場合があります。

### 4. サービスの確認

以下のURLでサービスにアクセスできます：

- **フロントエンド**: http://localhost:3000
- **バックエンドAPI**: http://localhost:8000
- **Hasura Console**: http://localhost:8080

### 5. データベースマイグレーションの適用

マイグレーションは自動的に適用されますが、手動で適用する場合は：

```bash
docker compose -f docker/docker-compose.yml exec hasura hasura-cli migrate apply --database-name default
docker compose -f docker/docker-compose.yml exec hasura hasura-cli metadata apply
```

## 使い方

### 初回ログイン

デフォルトの管理者アカウント：

- **メールアドレス**: admin@example.com
- **パスワード**: admin123

### 新規ユーザー登録

1. http://localhost:3000 にアクセス
2. 「こちら」リンクをクリックして登録ページへ
3. メールアドレスとパスワードを入力
4. 「登録」ボタンをクリック

### TODOの管理

1. ダッシュボードで新しいTODOを作成
2. タイトルと説明（任意）を入力
3. 完了したTODOにチェックを入れる
4. 編集・削除も可能

### 管理者機能

1. 管理者アカウントでログイン
2. 「管理画面」ボタンをクリック
3. ユーザー管理タブでユーザーの管理
4. 全TODOタブで全ユーザーのTODOを確認

## API エンドポイント

### 認証

- `POST /api/register` - ユーザー登録
- `POST /api/login` - ログイン
- `GET /api/profile` - プロフィール取得（要認証）

### TODO

- `GET /api/todos` - TODO一覧取得（要認証）
- `GET /api/todos/:id` - TODO詳細取得（要認証）
- `POST /api/todos` - TODO作成（要認証）
- `PUT /api/todos/:id` - TODO更新（要認証）
- `DELETE /api/todos/:id` - TODO削除（要認証）

### 管理者

- `GET /api/admin/users` - 全ユーザー取得（要管理者権限）
- `GET /api/admin/users/:id` - ユーザー詳細取得（要管理者権限）
- `PUT /api/admin/users/:id/role` - ユーザーロール変更（要管理者権限）
- `DELETE /api/admin/users/:id` - ユーザー削除（要管理者権限）
- `GET /api/admin/todos` - 全TODO取得（要管理者権限）

## 開発

### バックエンドの開発

```bash
cd backend
go mod download
go run cmd/server/main.go
```

### バックエンドのテスト実行方法

todo-backend-testコンテナで実行します：

```bash
docker-compose exec backend-test go mod download

# パッケージ全体のテストを詳細表示で実行
docker-compose exec backend-test go test ./internal/service -v

# 個別のテスト関数を実行（例）
docker-compose exec backend-test go test ./internal/service -run TestTodoService_GetTodos -v

# サブテストを指定して実行（名前にスペースがある場合はダブルクォート）
docker-compose exec backend-test go test ./internal/service -run "TestTodoService_UpdateTodo/with changes" -v

# race detector を有効にして実行
docker-compose exec backend-test go test ./internal/service -race -v
```

### フロントエンドの開発

```bash
cd frontend
npm install
npm run dev
```

### ログの確認

```bash
# 全サービスのログ
docker compose -f docker/docker-compose.yml logs -f

# 特定のサービスのログ
docker compose -f docker/docker-compose.yml logs -f backend
docker compose -f docker/docker-compose.yml logs -f frontend
docker compose -f docker/docker-compose.yml logs -f hasura
docker compose -f docker/docker-compose.yml logs -f postgres
```

## トラブルシューティング

### ポートが既に使用されている

既に使用されているポートがある場合、`docker/docker-compose.yml` のポート番号を変更してください。

### データベース接続エラー

PostgreSQLの起動を待ってから、バックエンドとHasuraが起動します。エラーが続く場合は：

```bash
docker compose -f docker/docker-compose.yml down -v
docker compose -f docker/docker-compose.yml up -d
```

### フロントエンドが起動しない

ビルドエラーの場合、依存関係を再インストール：

```bash
docker compose -f docker/docker-compose.yml down
docker compose -f docker/docker-compose.yml build --no-cache frontend
docker compose -f docker/docker-compose.yml up -d
```

## セキュリティに関する注意

本番環境にデプロイする場合は、必ず以下を変更してください：

1. `.env` ファイルの `JWT_SECRET` を強力なランダム文字列に変更
2. `HASURA_GRAPHQL_ADMIN_SECRET` を強力なパスワードに変更
3. データベースのパスワードを変更
4. デフォルトの管理者アカウントのパスワードを変更または削除
5. HTTPS を使用する
6. 適切な CORS 設定を行う

## ライセンス

MIT License

## サポート

問題が発生した場合は、GitHubのIssuesセクションで報告してください。
