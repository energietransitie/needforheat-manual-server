FROM golang:1.20 as build

WORKDIR /go/src/twomes-manual-server

# Create /source and /parsed folders to be copied later.
RUN mkdir /source && \
    mkdir /parsed

# Download dependencies.
COPY ./go.mod ./go.sum .
RUN go mod download

# Build healthcheck binary.
COPY ./cmd/healthcheck/ ./cmd/healthcheck/
RUN CGO_ENABLED=0 go build -o /go/bin/healthcheck ./cmd/healthcheck/

# Build server binary.
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/server ./cmd/server/

FROM gcr.io/distroless/static-debian11

# Copy /source and /parsed folders with correct permissions.
COPY --from=build --chown=nonroot /source /source
COPY --from=build --chown=nonroot /parsed /parsed

# Copy healthcheck binary.
COPY --from=build /go/bin/healthcheck /usr/bin/

# Copy server binary.
COPY --from=build /go/bin/server /

USER nonroot

VOLUME /source

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=30s --start-interval=2s --retries=3 \
    CMD ["healthcheck"]

CMD ["/server"]