IMAGE_NAME ?= ssh-proxy
pwd = `pwd`
DOCKER_PATH ?= /var/run/docker.sock
image:
	@docker build --rm -t ${IMAGE_NAME} .
ash:
	@docker run --rm -it -p 2222:2222  -v ${DOCKER_PATH}:${DOCKER_PATH} -v ${pwd}:/go/src/app -w /go/src/app ${IMAGE_NAME} ash

.PHONY: image
