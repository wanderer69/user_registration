cd ../migrations/init
goose postgres "postgres://user:password@localhost:5432/postgres?sslmode=disable" status
goose postgres "postgres://user:password@localhost:5432/postgres?sslmode=disable" up
cd ..
goose postgres "postgres://user:password@localhost:5432/user_registration?sslmode=disable" up
