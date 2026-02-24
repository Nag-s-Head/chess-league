FROM golang:latest AS build

RUN apt-get update
RUN apt-get install ca-certificates make npm -y
RUN apt-get upgrade -y
RUN npm i -g pnpm

WORKDIR /build

COPY . .
RUN go build

FROM build AS release

RUN useradd -m app
WORKDIR /home/app
USER app

EXPOSE 8080
COPY --from=build /build/chess-league .
CMD ["./chess-league"]
