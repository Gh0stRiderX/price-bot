FROM alpine:latest

COPY price-bot /usr/bin/price-bot

VOLUME /opt/price-bot/config /opt/config

ENTRYPOINT ["sh", "-c", "price-bot", "-smtp-filepath=/opt/config/smtp.json"]
