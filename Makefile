build:
	GOOS=js GOARCH=wasm go build -o ./web/main.wasm ./web/main.go