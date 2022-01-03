.PHONY: build
build: clear
	go build -o huawei-push-authorizator cmd/server/main.go 

.PHONY: docker-build
docker-build:
	docker build . -t huawei-push-authorizator
	
.PHONY: run
run: 
	HU_API_TOKEN= HU_CLIENT_ID= HU_CLIENT_SECRET= go run cmd/server/main.go
	
.PHONY: clear	
clear:
	rm huawei-push-authorizator || true


	