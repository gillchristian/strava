FROM node:20-alpine

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci --omit=dev

COPY server/ server/
COPY tsconfig.server.json ./

CMD ["npx", "tsx", "server/index.ts"]
