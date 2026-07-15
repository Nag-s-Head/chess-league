FROM golang:1.26.5-trixie AS initial

FROM initial AS with_go_mod
COPY ./go.mod .
RUN go mod download

FROM with_go_mod AS build
WORKDIR /build
RUN apt-get update && apt-get install -y npm
RUN npm install -g pnpm@10

COPY . .
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" --mount=type=cache,target="/build/node_modules" go generate ./... && go build

FROM initial AS release

RUN useradd -m app
WORKDIR /home/app
USER app

EXPOSE 8080
COPY --from=build /build/chess-league .
COPY ./knight.png .
CMD ["./chess-league"]
