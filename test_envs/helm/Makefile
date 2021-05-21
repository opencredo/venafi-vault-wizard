IMAGE_NAME=vault-with-venafi-plugins
PLUGIN_VERSION=0.9.0

.PHONY: build
build:
	@echo "Building Docker image with version ${PLUGIN_VERSION} of the Venafi plugins"
	docker build -t ${IMAGE_NAME} --build-arg PLUGIN_VERSION=${PLUGIN_VERSION} .
