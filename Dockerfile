FROM golang:1.15-alpine3.13 as build

WORKDIR /app

# Copy build context to WORKDIR
COPY . .

# Build GO binary
RUN go build -o ecr-refresh ./cmd/ecr-refresh

FROM alpine:3.13

COPY --from=build /app/ecr-refresh .

# TODO: Run as non-root user

CMD [ "./ecr-refresh"]
