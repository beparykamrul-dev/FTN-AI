FROM golang:1.24-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY services/control-plane/go.mod services/control-plane/go.sum ./
RUN GOPROXY=direct go mod download
COPY services/control-plane/*.go ./
COPY services/control-plane/internal ./internal
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/ftn-control-plane .

FROM gcr.io/distroless/static-debian12
COPY --from=build /out/ftn-control-plane /ftn-control-plane
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/ftn-control-plane"]
