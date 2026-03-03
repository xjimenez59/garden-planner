FROM golang AS debug-base

FROM debug-base as debug
WORKDIR /opt/app/meteo
COPY . .

# Download all the dependencies and install the package
RUN go get -d -v ./...
RUN go install -v ./...

# Build the Go app
#RUN go build -o ./bin/garden-planner-meteo
RUN go build -o ./bin/garden-planner-meteo -buildmode=exe
# Expose port to the outside world
EXPOSE 8082

# end debug
# =======================================================
# begin windows

FROM debug AS windows

# Run the executable
COPY --from=debug /opt/app/meteo/bin/garden-planner-meteo /usr/local/bin/
ENTRYPOINT ["garden-planner-meteo"]
#CMD ["garden-planner-meteo"]

# end windows
# =======================================================
# begin raspberry

# build an image on target arch
ARG BUILDPLATFORM=linux/arm64

FROM --platform=$BUILDPLATFORM golang:alpine AS raspberry-pi-build-base

FROM raspberry-pi-build-base AS raspberry-pi-build

WORKDIR /opt/app/meteo
COPY . .

# Download all the dependencies and install the package
RUN go get -d -v ./...
RUN go install -v ./...

# Build the Go app
ARG TARGETARCH=arm64
ARG TARGETOS=linux
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o ./bin/garden-planner-meteo

# Finally, copy the executable file on a fresh "minimal" distro
FROM --platform=$BUILDPLATFORM arm64v8/alpine AS raspberry-base

FROM raspberry-base AS raspberry 
COPY --from=raspberry-pi-build /opt/app/meteo/bin/garden-planner-meteo /usr/local/bin/


# Expose port to the outside world
EXPOSE 8082


# Run the executable
ENTRYPOINT ["garden-planner-meteo"]


# end raspberry
# =======================================================
# begin nas

# build an image on target arch
ARG BUILDPLATFORM=linux/amd64

FROM --platform=$BUILDPLATFORM golang:alpine AS nas-build-base

FROM nas-build-base AS nas-build
WORKDIR /opt/app/meteo
COPY . .

# Download all the dependencies and install the package
RUN go get -d -v ./...
RUN go install -v ./...

# Build the Go app
ARG TARGETARCH=amd64
ARG TARGETOS=linux
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o ./bin/garden-planner-meteo

# Finally, copy the executable file on a fresh "minimal" distro
FROM --platform=$BUILDPLATFORM alpine:latest AS nas-base

FROM nas-base AS nas 
COPY --from=nas-build /opt/app/meteo/bin/garden-planner-meteo /usr/local/bin/
# Expose port to the outside world
EXPOSE 8082


# Run the executable
ENTRYPOINT ["garden-planner-meteo"]




