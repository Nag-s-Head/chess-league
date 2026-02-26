FROM golang:latest AS initial

FROM initial AS with_go_mod
COPY ./go.mod .
RUN go mod download

FROM with_go_mod AS build
WORKDIR /build

COPY . .
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build

FROM initial AS release

RUN useradd -m app
WORKDIR /home/app
USER app

EXPOSE 8080
COPY --from=build /build/chess-league .
CMD ["./chess-league"]
