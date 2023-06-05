#Stage 1 - Install dependencies and build the app in a build environment
FROM debian:latest AS build-env

# Install flutter dependencies
RUN apt-get update
RUN apt-get install -y curl git wget unzip libgconf-2-4 gdb libstdc++6 libglu1-mesa fonts-droid-fallback lib32stdc++6 python3 sed
RUN apt-get clean

# Clone the flutter repo
RUN git clone https://github.com/flutter/flutter.git /usr/local/flutter

# Set flutter path
ENV PATH="${PATH}:/usr/local/flutter/bin:/usr/local/flutter/bin/cache/dart-sdk/bin"

# Run flutter doctor
RUN flutter doctor -v
RUN flutter channel master
RUN flutter upgrade

# Copy files to container and build
RUN mkdir /app/
COPY . /app/
WORKDIR /app/

RUN flutter build web

# Stage 2 - Create the run-time image
FROM nginx:1.21.1-alpine as debug
COPY --from=build-env /app/build/web /usr/share/nginx/html

# Expose port to the outside world
EXPOSE 80

# end debug
# ---------------------------------------------
# pour raspberry

# build an image on target arch
ARG BUILDPLATFORM=linux/arm64

#Stage 1 - Install dependencies and build the app in a build environment
FROM  --platform=$BUILDPLATFORM arm64v8/debian:latest AS build-rpi-base

# Install flutter dependencies
RUN apt-get update
RUN apt-get install -y curl git wget unzip libgconf-2-4 gdb lib32stdc++6-x32-cross libglu1-mesa fonts-droid-fallback python3 sed
RUN apt-get clean

# Clone the flutter repo
RUN git config --global pack.window 1
RUN git clone https://github.com/flutter/flutter.git /usr/local/flutter

# Set flutter path
ENV PATH="${PATH}:/usr/local/flutter/bin:/usr/local/flutter/bin/cache/dart-sdk/bin"

# Run flutter doctor
RUN flutter doctor -v
RUN flutter channel master
RUN flutter upgrade

FROM build-rpi-base as build-rpi

# Copy files to container and build
RUN mkdir /app/
COPY . /app/
WORKDIR /app/

RUN flutter build web

# Stage 2 - Create the run-time image
FROM  --platform=$BUILDPLATFORM arm64v8/nginx:1.21.1-alpine as raspberry
COPY --from=build-rpi /app/build/web /usr/share/nginx/html

# acces https. les certificats doivent etre deposes sur le serveur hôte dans un volume accessible
COPY ./nginx.conf /etc/nginx/conf.d/garden-planner.conf

# Expose port to the outside world
EXPOSE 80
