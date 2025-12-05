# dockerfile
FROM node:18-alpine AS base

# 依存関係インストール
FROM base AS deps
WORKDIR /app
COPY ./frontend/package.json ./
# package-lock.json や .npmrc がある場合は一緒にコピー
# COPY ./frontend/package-lock.json ./
# COPY ./frontend/.npmrc ./
RUN npm install

# ビルド
FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
# フロントエンドのソースのみコピー
COPY ./frontend/ .
# `public` が存在しない場合に備えて作成
RUN mkdir -p /app/public
ARG NEXT_PUBLIC_HASURA_ENDPOINT
ARG NEXT_PUBLIC_BACKEND_URL
ENV NEXT_PUBLIC_HASURA_ENDPOINT=$NEXT_PUBLIC_HASURA_ENDPOINT
ENV NEXT_PUBLIC_BACKEND_URL=$NEXT_PUBLIC_BACKEND_URL
RUN npm run build

# 実行用イメージ
FROM base AS runner
WORKDIR /app
ENV NODE_ENV production

RUN addgroup --system --gid 1001 nodejs \
 && adduser --system --uid 1001 nextjs

# `public` を配置（存在保証済み）
COPY --from=builder /app/public ./public
# Next.js の standalone 出力を配置（server.js を含む）
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs
EXPOSE 3000
ENV PORT 3000
CMD ["node", "server.js"]
