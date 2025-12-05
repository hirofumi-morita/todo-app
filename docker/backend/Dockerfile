# dockerfile
FROM golang:1.21-alpine AS builder

# バックエンドのモジュールルートを作業ディレクトリに設定
WORKDIR /app/backend

# modulesファイルのみコピー（go.sum は無くても可）
COPY backend/go.mod .
# go.mod だけでダウンロードを試行（足りない分は後で tidy）
RUN go mod download || true

# バックエンドのソースをコピー
COPY backend/ .

# 依存関係を解決して go.sum を生成
RUN go mod tidy

# バックエンドをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/server

# 実行ステージ
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8000
CMD ["./main"]
