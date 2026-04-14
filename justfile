[group('local')]
test:
  go test -coverprofile cover.out -coverpkg=./... -v ./...
  go tool cover -html cover.out -o cover.html

[group('local')]
make:
  go build

[group('local')]
lint:
  golangci-lint run --fix

[group('local')]
run release:
  docker run --rm -v "$PWD/_stor:/app/data" -v "$PWD/_build/:/app/data/plugins" -p 8080:80 gotify/server:{{release}}

[group('local')]
e2e release: (sync-versions "v" + release) (build "arm64") (run release)

plugin-name := "gotify-slack-webhook"
[group('CI')]
build arch os="linux" :
  #!/usr/bin/env bash
  set -euxo pipefail
  mkdir -p _build
  # NOTE: Drop 2 characs because toolchain is "go1.2.3" and image is tagged with just version
  version=$(go mod edit -json | jq -r '.Toolchain[2:]')
  docker run --rm -v "$PWD/.:/mnt" -w /mnt gotify/build:${version}-{{os}}-{{arch}} go build -mod=readonly -a -installsuffix cgo -ldflags="-w -s" -buildmode=plugin -o "_build/{{plugin-name}}-{{os}}-{{arch}}.so"

server-mod := "gotify-server.mod"
[group('CI')]
sync-versions release:
  # Downloads the target server versions dependency lockfile to sync it to the local one
  wget -LO {{server-mod}} https://raw.githubusercontent.com/gotify/server/{{release}}/go.mod
  echo "{{release}}" > SERVER_VERSION.txt
  # Syncs the downloaded server lockfile with the local one to make plugin binary compatible
  go run github.com/gotify/plugin-api/cmd/gomod-cap -from {{server-mod}} -to go.mod
  # Sync Toolchain and Go language versions as well
  go mod edit -go $(go mod edit -json gotify-server.mod | jq -r '.Go')
  go mod edit -toolchain $(go mod edit -json gotify-server.mod | jq -r '.Toolchain')
  go mod tidy

