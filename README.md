This is a Url Shortener app
How to run project:
git clone https://github.com/serhiirubets/go-url-shortener.git
Run docker-compose file: docker compose up -d
Install all dependencies go mod tidy
Run migration once go run migrations/auto.go
Run app go run cmd/main.go