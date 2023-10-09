# --------------------------------------------------------------------------------
# base image
# --------------------------------------------------------------------------------
FROM golang:1.21 AS build

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o nethack-notifier

# --------------------------------------------------------------------------------
# run image
# --------------------------------------------------------------------------------
FROM gcr.io/distroless/static

WORKDIR /app

COPY --from=build /app/nethack-notifier .

ENV RECORD_FILE_NAME="/home/nethack/nh367/record"
ENV IN_PROGRESS_DIR="/home/nethack/dgldir/inprogress-nh367"
ENV WEBHOOK_URL="https://sample.com"
ENV AVATAR_URL=""
ENV USER_NAME="Nethack Notifer"

CMD [ "/app/nethack-notifier" ]
