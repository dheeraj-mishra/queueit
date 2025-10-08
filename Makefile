build:
	@cd cmd/server; GOOS=linux GOARCH=amd64 go build -o ../../bin/queueit
	@cp cmd/server/.env ./bin
	@cd cmd/server; CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o ../../bin/queueit.exe
	@cp cmd/server/.env ./bin
	@echo "build success: path: ./bin/queueit"

run: build
	@cd bin; ./queueit

swagger-ui:
	@swag init
	@cp docs/swagger.json /home/dheeraj/swagger-setup/dist/.
	@echo "COPY URL: http://0.0.0.0:8081"
	@cd /home/dheeraj/swagger-setup/dist; python3 -m http.server 8081 > /dev/null 2>&1
