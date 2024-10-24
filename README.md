﻿# StreamSink.Go

 [![Build Status](https://teamcity.sedrad.com/app/rest/builds/buildType:(id:StreamSinkGo_Build)/statusIcon)](https://teamcity.sedrad.com/viewType.html?buildTypeId=StreamSinkGo_Build&guest=1)

StreamSink is an automated streaming recording server written in Go, which exposes a REST API.
It also manages the recorded catalogue, creates previews and has an API for cutting the recorded videos.

Which streaming services are supported? Just look at the `youtubedl` [youtube-dl supported sites](https://ytdl-org.github.io/youtube-dl/supportedsites.html)

The server will automatically keep checking streams and start recording them, once
they are online. When streams have finished or have been paused by the used, they will be added to the catalogue.

A Vue client [StreamSink.Vue](https://github.com/srad/StreamSink.Vue) is provided as a UI.

## Setup

You need to run the client and server separately in order to use the application.

### Server

The application is primarily crafted for Linux but also has code paths specifically for Windows but still does not work well on Windows,
due to process communication issues.
So, it is recommended to use Linux, since process management works on Linux.

The server requires that four applications being accessible from the environment: `ffmpeg`, `ffprobe`, `youtubedl`.
All are available for Linux and Windows. Below are more specific descriptions for a setup.

You can also use docker, but since `ffmpeg` is used with high CPU usage, you might want to run a native server instance.

#### Config

The server has one config file under `conf/app.yml` which must provide
the listed values. You can copy the default `conf.app.default.yml` to `conf.app.yml`.

### Client

A UI for the server is provided at, and can be deploy with the

### Linux

You can install all required applications via the default packet manager, but it is recommended to compile ffmpeg from source, since the
repo version of distros are typically far behind. However, you can still use the maintainer's version without any problem.

A setup script `setup.sh` is provided and installs all needed packages under Ubuntu and Debian and also compiles
ffmpeg from source.

### Windows

The server uses SQLite3 for storage and Go requires specific build tools in order to compile on Windows:

1. Install gcc from here http://tdm-gcc.tdragon.net/download
2. Install the font under `assets/..`
3. Download ffmpeg and place it into the global system path

#### Issue with Windows

The server needs read and execute permissions on processes on Windows.
It is also not enough to run the application as administrator, although
it will work without these permissions, Windows will cause an and error
exit code 255 (insufficient permission to manage process).

## Run Test

```go
go test ./...
```

## Notes & Limitations

1. All streaming services allow only a limited number of request made by each client.
If this limit is exceeded the client will be temporarily or permanently blocked.
In order to circumvent this issue, the application does strictly control the 
timing between each request. However, this might cause that the recording will only start
recording after a few minutes and not instantly.
2. The system has disaster recovery which means that if the system crashes during recordings,
it will try to recover all recordings on the next launch. However, due to the nature of
streaming videos and the crashing behavior, the video files might get corrupted.
In this case they will be automatically delete from the system, after they have been
checked for integrity. Otherwise, they are added to the library.
