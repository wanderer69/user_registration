generate:
	oapi-codegen -generate server -package public ./api/rest/swagger_public.yaml > ./pkg/api/public/server.go
	oapi-codegen -generate spec -package public ./api/rest/swagger_public.yaml > ./pkg/api/public/spec.go
	oapi-codegen -generate types -package public ./api/rest/swagger_public.yaml > ./pkg/api/public/types.go

up:
	docker compose -f ./deploy/docker-postgre-compose.yml up -d

down:
	docker compose -f ./deploy/docker-postgre-compose.yml down

run:
	sh ./start.sh

