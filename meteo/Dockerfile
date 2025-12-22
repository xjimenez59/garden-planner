FROM golang as debug-base

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

FROM debug as windows

# Run the executable
COPY --from=debug /opt/app/meteo/bin/garden-planner-meteo /usr/local/bin/
ENTRYPOINT ["garden-planner-meteo"]
#CMD ["garden-planner-meteo"]

# end windows
# =======================================================
# begin raspberry

# build an image on target arch
ARG BUILDPLATFORM=linux/arm64

FROM --platform=$BUILDPLATFORM golang:alpine as raspberry-pi-build-base

FROM raspberry-pi-build-base as raspberry-pi-build

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
FROM --platform=$BUILDPLATFORM arm64v8/alpine as raspberry-base

FROM raspberry-base as raspberry 
COPY --from=raspberry-pi-build /opt/app/meteo/bin/garden-planner-meteo /usr/local/bin/


# Expose port to the outside world
EXPOSE 8082


# Run the executable
ENTRYPOINT ["garden-planner-meteo"]



