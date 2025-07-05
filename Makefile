DIST := dist

build: build-web
	go build -o $(DIST)/gollery cmd/gollery/main.go


build-web:
	cd web \
    && rm -rf ../server/web \
    && npm run build \
    && mv dist ../server/web \
