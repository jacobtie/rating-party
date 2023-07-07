build:
	cd client && npm run build
	rm -rf server/cmd/dist
	mv client/dist server/cmd/
	cd server && go build -o app cmd/main.go
run:
	cd server && ./app
