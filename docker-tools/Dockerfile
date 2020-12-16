####### Start from the latest golang base image ###############
FROM golang:1.14.3-buster as builder
LABEL maintainer="Motiv Labs <dev@motivsolutions.com>"
WORKDIR /app
COPY ./ ./

# Build the Go app.  This below command builds Go statically linked so that all C libs etc... are included
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /usr/local/bin/main .

######## Start a new stage from a small image #######
FROM alpine:latest as prod
# Add root certs
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# Copy the Pre-built binary file from the builder stage
COPY --from=builder /usr/local/bin/main .
EXPOSE 8080
COPY docker-tools/entrypoint.sh /usr/local/bin
RUN chmod 755 /usr/local/bin/entrypoint.sh
# Add cqlsh to the image.
RUN apk add --update \
    bash \
    curl \
    py-pip
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
