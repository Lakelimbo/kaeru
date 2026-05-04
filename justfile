set dotenv-load
export CGO_ENABLED := "0"

frontend-generate := if path_exists("./web/build") == "false" { `go generate ./...` } else { `` }

serve:
    {{ frontend-generate }}
    go run ./cmd/kaeru serve

build:
    go generate ./...
    go build -trimpath -ldflags "-s -w" ./cmd/kaeru

clear:
    rm -rf ./web/build
    rm -rf .kaeru

clear-frontend:
    rm -rf ./web/build

test:
    go test ./... -v --cover

test-report:
    go-test ./... -v --cover -coverprofile=coverage.out
    go tool cover -html=coverage.out

format:
    go fmt ./...
    go run github.com/swaggo/swag/v2/cmd/swag@latest fmt -d ./

lint:
    golangci-lint run -c ./golangci.yml ./...

sql-migration name:
    go run github.com/pressly/goose/v3/cmd/goose@latest create {{ name }} sql -dir ./internal/database/migrations

sql-migrate-up:
    go run github.com/pressly/goose/v3/cmd/goose@latest sqlite ./.kaeru/db/kaeru.db up -dir ./internal/database/migrations

generate-openapi:
    go run github.com/swaggo/swag/v2/cmd/swag@latest init --dir ./internal/api -g routes.go -o ./internal/server/spec --parseFuncBody --parseDependency --v3.1
    pnpx openapi-typescript ./internal/server/spec/swagger.json -o ./web/src/lib/client/v1.d.ts --properties-required-by-default --root-types --root-types-no-schema-prefix
