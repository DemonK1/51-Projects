module crm_sqlite

go 1.21.6

// go get -u github.com/gofiber/fiber/v3
// Fiber 是一个受 Express 启发的 Web 框架
// 建立在 Fasthttp 之上，Fasthttp 是 Go 最快的 HTTP 引擎。旨在简化快速开发，同时考虑到零内存分配和性能

require (
	github.com/gofiber/fiber/v3 v3.0.0-20240124121856-755f133ac161
	gorm.io/driver/sqlite v1.5.4
	gorm.io/gorm v1.25.5
)

require (
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/gofiber/utils/v2 v2.0.0-beta.3 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.17.4 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-sqlite3 v1.14.19 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.51.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
)
