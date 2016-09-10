# imstor

> A simple app to have your **im**ages **stor**ed

This software is still under active development, and has not been vetted for production use cases.

imstor is an anonymous image sharing web app - images are uploaded with a POST request, are assigned
a URL and have thumbnails generated asynchonously.

### Currently available:
 - RESTful HTTP interface for uploading and retrieving images
 - Automatic async thumbnail creation, with configurable sizes
 - Disk-based (persistent) and in-memory storage
 - Easy to swap out backend storage engines and front-end request servers

### TODO
 - Configuration-file driven storage
 - Add additional fields for thumbnail sizing
 - Distributed storage (sharding)
 - Distributed thumbnail processing (use machinery?)


## Development
### Getting started
We assume you have a correctly configured `GOPATH` and Go development toolset

#### BE
 - `mkdir -p $GOPATH/src/github.com/denbeigh2000`
 - `git clone git@github.com:denbeigh2000/imstor.git imstor`
 - `cd imstor`

#### FE
 - `cd $GOPATH/src/github.com/denbeigh2000/imstor`
 - `git clone git@github.com:denbeigh2000/imstor-web.git web`
 - `git submodule update --init --recursive`
 - `cd web`

#### Running the server
Run these two commands in separate shells from the `imstor` directory:
 - Frontend: `cd web && npm run start`
 - Backend `go run cmd/http-server/main.go`

The frontend app will be accessible at `http://localhost:8080/`, and the API at
`http://localhost:8080/api/`.
