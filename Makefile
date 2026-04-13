.PHONY: build test lint docker up down clean seed demo release ui-build ui-dev ui-test allinone allinone-run allinone-stop

ui-build:
	cd ui && npm install && npm run build
	rm -rf internal/api/ui/dist
	mkdir -p internal/api/ui
	cp -r ui/dist internal/api/ui/dist

ui-dev:
	cd ui && npm run dev

ui-test:
	cd ui && npm test

build: ui-build
	go build -o bin/agenthound ./cmd/agenthound

test:
	go test ./... -v -race -count=1

lint:
	golangci-lint run ./...

docker:
	docker build -f docker/Dockerfile -t agenthound:dev .

up:
	docker compose -f docker/docker-compose.yml up -d

down:
	docker compose -f docker/docker-compose.yml down

clean:
	rm -rf bin/ coverage.out ui/dist internal/api/ui/dist

seed:
	@bash scripts/seed-test-data.sh

demo:
	@bash scripts/seed-demo.sh

release:
	goreleaser release --clean --snapshot

allinone:
	docker build -f docker/Dockerfile.allinone -t agenthound:allinone .

allinone-run:
	docker run -d --name agenthound -p 8080:8080 -v agenthound-data:/data --restart unless-stopped agenthound:allinone

allinone-stop:
	docker stop agenthound && docker rm agenthound
