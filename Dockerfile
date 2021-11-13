# Building the binary of the App
FROM golang:1.15 AS build

# `goth` should be replaced with your project name
WORKDIR /go/src/goth

# Copy dependency files
COPY go.mod .
COPY go.sum .

# Downloads all the dependencies in advance (could be left out, but it's more clear this way)
RUN go mod download

# Copy all the Code and stuff to compile everything
COPY . .

# Builds the application as a staticly linked one, to allow it to run on alpine
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Moving the binary to the 'final Image' to make it smaller
FROM alpine:latest

WORKDIR /app

# Create the `public` dir and copy all the assets into it
#RUN mkdir ./static
#COPY ./static ./static

# `goth` should be replaced here as well
COPY --from=build /go/src/goth/app .

# Exposes port 3333 because our program listens on that port
EXPOSE 3333

CMD ["./app"]
