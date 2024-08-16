FROM node:22.6.0 as node

WORKDIR /app/ui
COPY ui/package-lock.json ./package-lock.json
COPY ui/package.json ./package.json
RUN npm ci

COPY ui /app/ui
RUN npm run build

FROM golang:1.23 as api

WORKDIR /app
COPY api/go.mod api/go.sum ./
RUN go mod download
COPY api /app
RUN CGO_ENABLED=0 go build

FROM alpine:3.7

COPY --from=api /app/api /api
COPY --from=node /app/ui/build /ui/build

RUN chmod +x /api

CMD ["/api"]