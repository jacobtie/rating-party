build:
	cd client && npm run build
	rm -rf server/cmd/dist
	mv client/dist server/cmd/

run:
	cd server && go run ./...
