ARG quay_expiration=never
ARG release_tag=0.0.0

# Build the manager binary
FROM golang:1.24 AS builder
ARG TARGETOS
ARG TARGETARCH
ARG release_tag

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the code
COPY cmd/main.go cmd/main.go
COPY Makefile Makefile
COPY hack/ hack/
COPY api/ api/
COPY internal/ internal/

# Copy git repo for sha info
COPY .git .git

RUN make build VERSION=${release_tag}

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/bin/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]
