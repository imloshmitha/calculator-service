FROM golang:1.22-alpine AS build

WORKDIR /src
COPY backend/go.mod ./
COPY backend/cmd ./cmd
COPY backend/internal ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/server ./cmd/server

FROM alpine:3.21

COPY --from=build /out/server /server
ENV PORT=8080
EXPOSE 8080
USER nobody
ENTRYPOINT ["/server"]
