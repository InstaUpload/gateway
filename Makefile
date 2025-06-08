run: docs
	go run *.go

docs:
	swag fmt
	swag init
