DIST := dist

build: build-web
	go build -o $(DIST)/gollery cmd/gollery/main.go


build-web:
	cd web \
    && rm -rf ../server/web \
    && npm install \
    && npm run build \
    && mv dist ../server/web \
