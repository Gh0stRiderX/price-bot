FROM arm32v7/alpine:latest

COPY price-bot /usr/bin/price-bot

EXPOSE 8091

ENTRYPOINT ["sh", "-c", "price-bot"]
