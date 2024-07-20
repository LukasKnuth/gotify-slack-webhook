server-mod := "gotify-server.mod"
server-go-version := "1.22.4"
plugin-name := "gotify-slack-webhook"

build arch:
  mkdir -p _build
  docker run --rm -v "$PWD/.:/build" -w /build "gotify/build:{{server-go-version}}-{{arch}}" go build -mod=readonly -a -installsuffix cgo -ldflags "$$LD_FLAGS" -buildmode=plugin -o _build/{{plugin-name}}-{{arch}}.so

build-all: (build "linux-amd64") (build "linux-arm64")

download-gotify-server-mod:
  wget -LO {{server-mod}} https://raw.githubusercontent.com/gotify/server/master/go.mod
  echo "Also note that the Go version must match with Gotiy server (set via ASDF)"

verify-versions: download-gotify-server-mod
   go run github.com/gotify/plugin-api/cmd/gomod-cap -from {{server-mod}} -to go.mod
   go mod tidy
