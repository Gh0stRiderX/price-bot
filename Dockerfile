FROM arm32v7/alpine:latest

COPY price-bot /usr/bin/price-bot

VOLUME /opt/price-bot/config /opt/config

EXPOSE 8091

ENTRYPOINT ["sh", "-c", "price-bot", "-smtp-filepath=/opt/config/smtp.json -port=8091"]
