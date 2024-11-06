//go:build go1.23

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=internal/oapi/cfg.yaml internal/oapi/tltv.yaml
//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc generate -f db/sqlc.yaml
//go:generate go run github.com/golang/mock/mockgen -package mockdb -destination internal/mock/db/store.go talkliketv.click/tltv/db/sqlc Querier
//go:generate go run github.com/golang/mock/mockgen -package mockt -destination=internal/mock/translates/translates.go -source=internal/translates/translates.go
//go:generate go run github.com/golang/mock/mockgen -package mockc -destination=internal/mock/clients/clients.go -source=internal/translates/clients.go
//go:generate go run github.com/golang/mock/mockgen -package mocka -destination=internal/mock/audiofile/audiofile.go -source=internal/audio/audiofile/audiofile.go

package main
