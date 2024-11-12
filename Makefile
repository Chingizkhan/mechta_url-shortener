swagger:
	swagger generate spec -w cmd/app -o ./swagger.yaml -m

app:
	go run cmd/app/app.go
