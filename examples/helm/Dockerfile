FROM vault:1.7.1

ARG PKI_MONITOR_PLUGIN_VERSION=0.9.0
ARG PKI_BACKEND_PLUGIN_VERSION=0.9.0

RUN set -eux; \
    apkArch="$(apk --print-arch)"; \
    case "$apkArch" in \
        x86_64) ARCH='' ;; \
        x86) ARCH='86' ;; \
        *) echo >&2 "error: unsupported architecture: $apkArch"; exit 1 ;; \
    esac && \
    mkdir -p /tmp/plugin_dl && \
    cd /tmp/plugin_dl && \
    wget https://github.com/Venafi/vault-pki-monitor-venafi/releases/download/v${PKI_MONITOR_PLUGIN_VERSION}/venafi-pki-monitor_v${PKI_MONITOR_PLUGIN_VERSION}_linux${ARCH}_optional.zip && \
    wget https://github.com/Venafi/vault-pki-backend-venafi/releases/download/v${PKI_BACKEND_PLUGIN_VERSION}/venafi-pki-backend_v${PKI_BACKEND_PLUGIN_VERSION}_linux${ARCH}.zip && \
    mkdir -p /vault/plugins && \
    unzip venafi-pki-monitor_v${PKI_MONITOR_PLUGIN_VERSION}_linux${ARCH}_optional.zip && \
    cp venafi-pki-monitor_optional /vault/plugins/venafi-pki-monitor_v${PKI_MONITOR_PLUGIN_VERSION} && \
    unzip venafi-pki-backend_v${PKI_BACKEND_PLUGIN_VERSION}_linux${ARCH}.zip && \
    cp venafi-pki-backend /vault/plugins/venafi-pki-backend_v${PKI_BACKEND_PLUGIN_VERSION} && \
    cd /tmp && \
    rm -rf /tmp/plugin_dl
