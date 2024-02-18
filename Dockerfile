FROM golang:1.22.0-bookworm

RUN apt update && apt upgrade -y && apt install unzip wget -y

RUN apt install chromium chromium-driver -y

WORKDIR /app

COPY . .

RUN go build .

ARG QUOTE_URL="https://www.yourquote.in/PROFILE_ID/quotes"
ENV QUOTE_URL=${QUOTE_URL}

# ENTRYPOINT [ "sleep","infinity" ]
ENTRYPOINT [ "./yourquote" ]