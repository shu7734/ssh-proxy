IMAGE_NAME ?= ssh-proxy
pwd = `pwd`
image:
	@docker build --rm -t ${IMAGE_NAME} .
ash:
	@docker run --rm -it -v ${pwd}:/go/src/app -w /go/src/app ${IMAGE_NAME} ash

.PHONY: image
