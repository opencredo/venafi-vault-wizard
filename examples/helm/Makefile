IMAGE_NAME=vault-with-venafi-plugins
PKI_MONITOR_PLUGIN_VERSION=0.9.0
PKI_BACKEND_PLUGIN_VERSION=0.10.3

.PHONY: build
build:
	@echo "Building Docker image with version ${PKI_MONITOR_PLUGIN_VERSION} of the pki monitor and version ${PKI_BACKEND_PLUGIN_VERSION} of the pki backend venafi plugins"
	docker build -t ${IMAGE_NAME} --build-arg PKI_MONITOR_PLUGIN_VERSION=${PKI_MONITOR_PLUGIN_VERSION} --build-arg PKI_BACKEND_PLUGIN_VERSION=${PKI_BACKEND_PLUGIN_VERSION} .
