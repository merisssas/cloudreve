# syntax=docker/dockerfile:1

FROM alpine:3.20

WORKDIR /cloudreve

# Timezone Frankfurt (Germany) -> Europe/Berlin
# Bisa override saat run: -e TZ=Europe/Berlin
ENV TZ=Europe/Berlin

RUN set -eux; \
    apk add --no-cache \
      ca-certificates \
      tzdata \
      vips-tools \
      ffmpeg \
      libreoffice \
      aria2 \
      supervisor \
      font-noto \
      font-noto-cjk \
      libheif \
      libraw-tools \
    ; \
    update-ca-certificates; \
    ln -snf "/usr/share/zoneinfo/${TZ}" /etc/localtime; \
    echo "${TZ}" > /etc/timezone; \
    mkdir -p /cloudreve/data/temp/aria2; \
    chmod -R 775 /cloudreve/data/temp/aria2

# Cloudreve feature flags (tetap seperti punyamu)
ENV CR_ENABLE_ARIA2=1 \
    CR_SETTING_DEFAULT_thumb_ffmpeg_enabled=1 \
    CR_SETTING_DEFAULT_thumb_vips_enabled=1 \
    CR_SETTING_DEFAULT_thumb_libreoffice_enabled=1 \
    CR_SETTING_DEFAULT_media_meta_ffprobe=1 \
    CR_SETTING_DEFAULT_thumb_libraw_enabled=1

# Copy configs + entrypoint
COPY .build/aria2.supervisor.conf /cloudreve/aria2.supervisor.conf
COPY .build/entrypoint.sh /cloudreve/entrypoint.sh

# Copy binary Cloudreve (WAJIB ARM64 kalau buat Oracle ARM)
COPY cloudreve /cloudreve/cloudreve

RUN set -eux; \
    chmod 755 /cloudreve/cloudreve /cloudreve/entrypoint.sh

EXPOSE 5212 443

VOLUME ["/cloudreve/data"]

ENTRYPOINT ["/bin/sh", "/cloudreve/entrypoint.sh"]
