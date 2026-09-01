FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci --quiet
COPY . .
RUN npm run build --if-present

FROM node:20-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production --quiet
COPY --from=builder /app/dist ./dist
EXPOSE 3000
CMD ["npm", "run", "dev"]
