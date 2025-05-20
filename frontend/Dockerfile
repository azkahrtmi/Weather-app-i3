# Build Stage
FROM node:22-alpine AS builder
WORKDIR /app
COPY . .
RUN npm ci
RUN npm run build
CMD [ "npm", "run", "dev", "--", "--host", "0.0.0.0" ]
