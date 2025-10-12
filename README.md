## Matchee Services - Go Gin Base (Clean Architecture)

### Stack
- Gin, GORM (MySQL)
- Redis, RabbitMQ
- golang-migrate (migrations)

### Structure
```
cmd/api/main.go
internal/
  config/
  entity/
  repository/
  usecase/
  controller/
  router/
pkg/
  database/
  redis/
  rabbitmq/
  logger/
migrations/
```

### Setup
- Tạo file `.env` từ `.env.example` (hoặc export biến môi trường tương ứng)
- Tạo DB MySQL `matchee` và cập nhật `MYSQL_DSN` nếu cần

### Commands
```bash
make tidy        # đồng bộ module
make migrate-up  # chạy migrations (yêu cầu cài migrate CLI)
make run         # chạy server :8080
```

### Endpoints
- `GET /health`
- `POST /api/users` body: `{ "email": "a@b.c", "name": "A" }`
- `GET /api/users/:id`


