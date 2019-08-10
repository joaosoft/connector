fmt:
	go fmt ./...

vet:
	go vet ./*

gometalinter:
	gometalinter ./*

dependency-get:
	./bin/dependency get

dependency-update:
	./bin/dependency update

dependency-reset:
	./bin/dependency reset