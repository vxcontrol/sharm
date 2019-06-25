# Sharm Agent

[![Build Status](https://travis-ci.org/vxcontrol/sharm.svg?branch=master)](https://travis-ci.org/vxcontrol/sharm)

Sharm is a public service to recording you terminal and to showing to everybody as real time broadcasting. Also, you are able to record and broadcast your voice and to attach chat comments to final result and to share this.

### Sharm key features:

  - Totally free service
  - Broadcast terminal and sound
  - Chat in real time broadcast
  - Supports MacOS, Linux and Windows
  - Unlimited storage of records

### Support

It's main page of the service and you should post your questions, issues and suggestions to this tracker page.
While the service is still young, we will welcome your suggestions for its improvement.

**Important**: We aren't responsible for any risks and losses associated with using of this service and this code.

## Dependencies

  - Windows version of the agent based on [winpty project](https://github.com/rprichard/winpty)
  - Transport subsystem has based on [websockets](https://github.com/gorilla/websocket) and TLS connection over golang library
  - Logging subsystem has based on [logrus](https://github.com/sirupsen/logrus)

## Building

It's very simple and native procedure:

```sh
$ cd sharm
$ go get
$ go build -o build/sharm
```

Of course you can use **GOARCH** and **GOOS** environment variables on build.

## Using

The agent working with original [Sharm service](https://sharm.io) and to use it you need to follow the instructions from this service.

## Changelog

### sharm v1.0

  - Supports MacOS, Linux and Windows
  - Share one screen at one time
  - Control from environment variables
  - Protocol v1 over JSON struct based on WebSocket

## Copyright

This project is distributed under the BSD 3-Clause license (see the LICENSE file in the project root).

By submitting a pull request for this project, you agree to license your contribution under the BSD 3-Clause license to this project.

© 2019  [vxcontrol.com](https://vxcontrol.com)™  [sharm.io](https://sharm.io)™
