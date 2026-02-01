# nethack-notifier

Application to notify nethack server activity via discord webhook.

[![Docker Hub](https://img.shields.io/badge/Docker%20Hub-Repository-2496ED?logo=docker&logoColor=white)](https://hub.docker.com/r/4nm1tsu/nethack-notifier)
[![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/4nm1tsu/nethack-notifier/docker-build.yml)](https://github.com/4nm1tsu/nethack-notifier/actions)



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
