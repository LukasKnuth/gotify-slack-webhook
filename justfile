server-mod := "gotify-server.mod"

build:
  go build -buildmode=plugin

download-gotify-server-mod:
  wget -LO {{server-mod}} https://raw.githubusercontent.com/gotify/server/master/go.mod
  echo "Also note that the Go version must match with Gotiy server (set via ASDF)"

verify-versions: download-gotify-server-mod
   go run github.com/gotify/plugin-api/cmd/gomod-cap -from {{server-mod}} -to go.mod
   go mod tidy
