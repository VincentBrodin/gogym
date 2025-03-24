build:
	$(MAKE) build-backend
	$(MAKE) build-frontend

build-backend:
	$(MAKE) -C backend build

build-frontend:
	$(MAKE) -C frontend build

build-docker:
	$(MAKE) -C frontend build
	docker build -f Dockerfile -t vincentbrodin/gogym:latest --no-cache .

push-docker:
	$(MAKE) build-docker
	docker push vincentbrodin/gogym:latest

start-docker:
	docker run -d -p 5050:8080 --name gogym vincentbrodin/gogym:latest

up-docker:
	-docker stop gogym
	-docker rm gogym
	$(MAKE) build-docker
	$(MAKE) start-docker

dev:
	$(MAKE) -C frontend dev 
	$(MAKE) -C backend dev 
