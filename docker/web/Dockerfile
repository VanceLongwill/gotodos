FROM golang

# Set shell to bash
# SHELL ["/bin/bash", "-c"]

# Install tools required for project
# Run `docker build --no-cache .` to update dependencies
# RUN apk add --no-cache git
RUN go get github.com/golang/dep/cmd/dep

# List project dependencies with Gopkg.toml and Gopkg.lock
# These layers are only re-built when Gopkg files are updated
COPY Gopkg.lock Gopkg.toml /go/src/github.com/vancelongwill/gotodos/
WORKDIR /go/src/github.com/vancelongwill/gotodos
# Install library dependencies
RUN dep ensure -vendor-only

# Copy the entire project and build it
# This layer is rebuilt when a file changes in the project directory
COPY . /go/src/github.com/vancelongwill/gotodos
RUN go build -o /bin/gotodos

# This results in a single layer image
#  FROM scratch
#  COPY --from=build /bin/gotodos /bin/gotodos
ENTRYPOINT ["/bin/gotodos"]
EXPOSE 80
CMD ["--help"]
