[group('local')]
test:
  go test ./...

[group('local')]
make:
  go build

server-mod := "gotify-server.mod"
docker-image := "gotify-build-arm64"
server-go-version := "1.22.4"
plugin-name := "gotify-slack-webhook"

[group('gotify')]
build-image:
  # This is here because the official gotify/builder images are AMD64 only and I'm on a M1 Mac Book
  docker build . -f Dockerfile.build -t {{docker-image}} --build-arg GO_VERSION={{server-go-version}}

[group('gotify')]
_build arch:
  mkdir -p _build
  docker run --rm -v "$PWD/.:/build" -w /build {{docker-image}} go build -mod=readonly -a -installsuffix cgo -ldflags="-w -s" -buildmode=plugin -o _build/{{plugin-name}}-{{arch}}.so

[group('gotify')]
build: (_build "linux-arm64")

[group('gotify')]
run:
  docker run --rm -v "$PWD/_build:/app/data/plugins" -p 8080:80 gotify/server-arm64

[group('gotify')]
download-gotify-server-mod:
  wget -LO {{server-mod}} https://raw.githubusercontent.com/gotify/server/master/go.mod
  echo "Also note that the Go version must match with Gotiy server (set via ASDF)"

[group('gotify')]
verify-versions: download-gotify-server-mod
   go run github.com/gotify/plugin-api/cmd/gomod-cap -from {{server-mod}} -to go.mod
   go mod tidy
