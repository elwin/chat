# Websocket chat
To play around a bit with websockets, I built this little chat app. It is powered by Go on the backend and React on the
frontend. 

# Build
## Prerequisites
- Go Toolchain
- Node.js
- Docker

## Actual Build
Check out the Makefile on how to build it. The project is aimed to be run on chat.elwin.dev with some
hardcoded values. To run it locally, you need to replace the base url in `ui/src/index.js` with `localhost:8888`.
After that, run
- `make build`
- `docker run -p 8888:8888 chat`