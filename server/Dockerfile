FROM golang:latest AS builder
WORKDIR /server/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM scratch
WORKDIR /bin/
COPY --from=builder /server/ .
CMD [ "./app" ]
EXPOSE 8000


