# generate public
oapi-codegen -generate server -package public ../api/rest/swagger_public.yaml > ../pkg/api/public/server.go
oapi-codegen -generate spec -package public ../api/rest/swagger_public.yaml > ../pkg/api/public/spec.go
oapi-codegen -generate types -package public ../api/rest/swagger_public.yaml > ../pkg/api/public/types.go
