//go:build go1.23

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=internal/oapi/cfg.yaml internal/oapi/tltv.yaml
//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc generate -f db/sqlc.yaml
//go:generate go run github.com/golang/mock/mockgen -package mockdb -destination db/mock/store.go talkliketv.click/tltv/db/sqlc Querier
//go:generate go run github.com/golang/mock/mockgen -destination=internal/mock/translates.go -source=api/translates.go

package main
