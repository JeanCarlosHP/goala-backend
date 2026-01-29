### **Backend**
* Golang + Fiber
* pgxpool + SQLC + PostgreSQL
* go-migrate v4
* Firebase Auth (tokens)
* AWS S3 (fotos)
* Prometheus + Grafana (monitoramento)
* Slog (logs estruturados com a interface domain.Logger)
* Viper (configuração)
* Clean Architecture + CQRS leve

- Caso seja nescessario utilize o MCP #io.github.upstash/context7 para consultar documentação técnica.
- Não adicione comentários em excesso.
- Não crie documentação extra.
- Utilize `go vet ./...` e `staticcheck ./...` para garantir a qualidade do código.