REGISTRY=docker.pkg.github.com/elwin
VERSION=latest

build-ui:
	cd ui; npm run build;

build-api:
	cd api; GOOS=linux GOARCH=amd64 go build;

build: build-api build-ui
	docker build . -t chat

publish: build
	docker tag chat $(REGISTRY)/chat/chat:$(VERSION)
	docker push $(REGISTRY)/chat/chat:$(VERSION)
