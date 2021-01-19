IMAGE_PREFIX ?= $(USER)/
IMAGE_TAG ?= latest
BUILD_BASE ?= golang:1.15-buster
DEPLOY_BASE ?= gcr.io/distroless/static:latest

.PHONY: build-images
build-images:
	docker pull ${BUILD_BASE}
	docker pull ${DEPLOY_BASE}
	docker build --build-arg BUILD_BASE=${BUILD_BASE} --build-arg DEPLOY_BASE=${DEPLOY_BASE} -t kaml:$(IMAGE_TAG) -f Dockerfile.kaml .
	docker build --build-arg BASE_IMAGE=kaml:$(IMAGE_TAG) -t kaml-set-annotations:$(IMAGE_TAG) -f Dockerfile.kaml-set-annotations .

.PHONY: build-images
push-images:
	docker tag kaml:$(IMAGE_TAG) $(IMAGE_PREFIX)kaml:$(IMAGE_TAG)
	docker push $(IMAGE_PREFIX)kaml:$(IMAGE_TAG)
	docker tag kaml-set-annotations:$(IMAGE_TAG) $(IMAGE_PREFIX)kaml-set-annotations:$(IMAGE_TAG)
	docker push $(IMAGE_PREFIX)kaml-set-annotations:$(IMAGE_TAG)
