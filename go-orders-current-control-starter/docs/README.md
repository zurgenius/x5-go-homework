# Swagger

Документация генерируется из аннотаций обработчиков:

```bash
go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g cmd/api/main.go -o docs
```

После запуска сервера UI доступен на `/swagger/index.html`.
