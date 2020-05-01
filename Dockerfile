# GitHub:       https://github.com/nicksherron
FROM golang:1.13-alpine AS build

ARG VERSION
ARG GIT_COMMIT
ARG BUILD_DATE

ARG CGO=1
ENV CGO_ENABLED=${CGO}
ENV GOOS=linux
ENV GO111MODULE=on

WORKDIR /go/src/github.com/nicksherron/proxi

COPY . /go/src/github.com/nicksherron/proxi/

# gcc/g++ are required to build SASS libraries for extended version
RUN apk update && \
    apk add --no-cache gcc g++ musl-dev


RUN go build -ldflags "-X github.com/nicksherron/proxi/internal.Version=${VERSION} -X github.com/nicksherron/proxi/cmd.Build=${GIT_COMMIT} -X github.com/nicksherron/proxi/cmd.BuildDate=${BUILD_DATE}  -extldflags=-static -extldflags=-lm" -o /go/bin/proxi

# ---

FROM alpine:3.11

COPY --from=build /go/bin/proxi /usr/bin/proxi

# libc6-compat & libstdc++ are required for extended SASS libraries
# ca-certificates are required to fetch outside resources (like Twitter oEmbeds)
RUN apk update && \
    apk add --no-cache ca-certificates libc6-compat libstdc++

VOLUME /data
WORKDIR /data

# Expose port for live server
EXPOSE 4444

CMD ["/usr/bin/proxi", "server", "--init"]
