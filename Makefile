REGISTRY="registry.digitalocean.com/studentelwin"

build-ui:
	cd ui; npm run build;

build-api:
	cd api; GOOS=linux GOARCH=amd64 go build;

publish: build-api build-ui
	docker build . -t chat
	docker tag chat $(REGISTRY)/chat
	docker push $(REGISTRY)/chat