FROM golang:1.20rc3-alpine3.17 AS build

RUN apk add --no-cache git
WORKDIR /src
ADD ./ /src
RUN CGO_ENABLED=0 go build -installsuffix 'static'

FROM gcr.io/distroless/static

COPY --from=build --chown=nonroot:nonroot /src/pub /root/pub
WORKDIR /root/
EXPOSE 9999
ENTRYPOINT ["/root/pub"]

CMD ["--driver", "sqlite", "--dsn" ,"/data/db.sqlite", "serve"]
