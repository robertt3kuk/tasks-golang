FROM golang:latest AS builder
WORKDIR /salter/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM scratch
WORKDIR /bin/
COPY --from=builder /salter/app .
CMD [ "./app" ]
EXPOSE 3000
