# Moneybots Tick Data Service

## Project Structure

```text
moneybotstds/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   │   └── publish_handler.go
│   │   ├── middleware/
│   │   │   └── auth.go
│   │   │   └── logger.go
│   │   └── routes.go
│   ├── config/
│   │   └── config.go
│   ├── models/
│   │   └── models.go
│   ├── repository/
│   │   └── db.go
│   │   └── redis.go
│   │   └── repository.go
│   └── service/
│       └── db_service.go
│       └── ticker_service.go
├── pkg/
│   └── response/
│       └── response.go
├── build/
│   └── bin/
│       └── mbtickservice
├── docs/
│   └── api.md
├── go.mod
├── go.sum
└── README.md
```
