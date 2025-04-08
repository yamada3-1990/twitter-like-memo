# フロントエンド用のDockerfile

FROM node:22.12.0-alpine
WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install --frozen-lockfile

# アプリを動かすために必要なフォルダをコピー
COPY . .

# コンテナのポートを公開
EXPOSE 5173

CMD ["npm", "run", "dev"]