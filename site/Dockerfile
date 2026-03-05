# do this in a sepearte container to cache the npm install step between rebuilds
FROM node:25-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
ENV SHOULD_USE_ADAPTER_NODE=1
RUN npm run build
RUN npm prune --production

FROM node:25-alpine
WORKDIR /app
COPY --from=builder /app/build build
COPY --from=builder /app/node_modules node_modules
COPY package.json .
EXPOSE 3000
ENV SHOULD_USE_ADAPTER_NODE=1
ENV NODE_ENV=production
CMD ["node", "build"]
