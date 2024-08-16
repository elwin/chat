REGISTRY=docker.pkg.github.com/elwin
VERSION=latest

build:
	docker build . -t chat

publish: build
	docker tag chat $(REGISTRY)/chat/chat:$(VERSION)
	docker push $(REGISTRY)/chat/chat:$(VERSION)

run: build
	docker run -p 8888:8888 chat