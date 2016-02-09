FROM golang:latest
RUN bash -c "go get github.com/scottjab/catbott"
# Use exposed env vars to configure service.
CMD bash -c "catbot /go/src/github.com/scottjab/catbot/config.json.example"
