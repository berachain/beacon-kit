module github.com/berachain/beacon-kit/mod/api

go 1.22.2

require (
	github.com/berachain/beacon-kit/mod/api/handlers v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi/v5 v5.0.12
	github.com/go-chi/cors v1.2.1
)

replace github.com/berachain/beacon-kit/mod/api/handlers => ./handlers
