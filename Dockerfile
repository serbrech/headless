FROM golang as build
WORKDIR /src
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o broadcast .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=build /src/broadcast ./
RUN chmod +x ./broadcast
ENTRYPOINT ["./broadcast"]