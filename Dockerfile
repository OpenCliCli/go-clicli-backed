FROM golang:alpine AS development
WORKDIR /app
ENV ENV=development
COPY . /app
RUN go get -v
RUN go get github.com/pilu/fresh
ENTRYPOINT ["fresh"]
EXPOSE 8080 3306

# FROM alpine:latest AS production
# WORKDIR /app
# RUN go build -o app
# # COPY --from=development /app .
# EXPOSE 8080
# ENTRYPOINT ["./app"]