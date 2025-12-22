reference: https://platform.openai.com/docs/guides/structured-outputs?example=chain-of-thought

also don't forget to add CI/CD

all required
put null for optional fields


test commands:
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html