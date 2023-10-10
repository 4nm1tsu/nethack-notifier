# nethack-notifier

<p align="center">
  <p align="center">Application to notify nethack server activity via discord webhook</p>
  <p align="center"><a href="https://github.com/4nm1tsu/nethack-notifier/actions/workflows/docker-build.yml"><img src="https://github.com/4nm1tsu/nethack-notifier/actions/workflows/docker-build.yml/badge.svg"></img></a></p>
  <p align="center">
    <a href="https://hub.docker.com/r/4nm1tsu/nethack-notifier/tags"><img alt="Docker Hub" src="http://dockeri.co/image/4nm1tsu/nethack-notifier"></a>
  </p>
</p>

## Usage

Run as a [nethack-server](https://github.com/4nm1tsu/nethack-server) sidecar container using docker-compose or on kubernetes.

### Environment Variables

| key              | type   | description                                  |
|------------------|--------|----------------------------------------------|
| RECORD_FILE_NAME | string | Path of record file.                         |
| IN_PROGRESS_DIR  | string | Path of inprogress directory of dgamelaunch. |
| WEBHOOK_URL      | string | Path of webhook URL of discord.              |
| AVATAR_URL       | string | Path of bot's avatar image.                  |
| USER_NAME        | string | Bot's name.                                  |
| SERVER_DOMAIN    | string | Domain name of the server                    |
