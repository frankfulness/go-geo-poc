# Go geolocation POC

1. `go version`, using `1.23.6` for this proof of concept

- If different version, see <https://go.dev/doc/manage-install> to `go install golang.org/dl/go1.23.6@latest`
- then download `go1.23.6 download`
- `go1.23.6 env GOROOT` should show version downloaded, however if we `go version` it may not be correct. Check `which go` to determine.
- Need be update your `.bashrc` or `.zshrc` and `source ~/.zshrc` after the changes:

    ```
    export GOROOT=/Users/zion/sdk/go1.23.6
    export GOPATH=$HOME/go # in home dir create dir called go and where everything is gonna go in there.
    # export PATH=$PATH:$GOPATH/bin # Adding to my path and path is stuff I can run from my term and concatinating in go/bin all in my path to be useable
    export PATH=$GOROOT/bin:$PATH:$GOPATH/bin  # Modified to include GOROOT/bin first
    ```

- Now `which go` and `go version` should reflect 1.23.6 which we are utilizing, though feel free to attempt with new Go versions.
- To summarize: The PATH setup export PATH=$GOROOT/bin:$PATH:$GOPATH/bin maintains all the functionality you currently have, while just ensuring you're using the new Go version. Here's what each part does:

```
$GOROOT/bin - This adds your new Go binary location
$PATH - This keeps all your existing PATH entries
$GOPATH/bin - This maintains your current setup that lets you run Go programs from anywhere

The order matters here - we put $GOROOT/bin first so it takes precedence for the go command itself, but everything else in your PATH (including your GOPATH/bin) stays exactly the same. This means:

You can still run go commands from any directory
Your installed Go programs will still be accessible from anywhere
Your existing Go projects and workspace structure stay the same
```

The only thing that changes is which version of Go you're using - everything else about your setup remains functional just as it is now.

2. Obtain gorilla mux for routing, uber h3 for geospatial indexing, & spherand for generating random points on a sphere.

```
go get -u github.com/gorilla/mux
go get -u github.com/uber/h3-go
go get -u github.com/mmcloughlin/spherand
```

3. After completing this POC, demo by `go run main.go` and navigate to `http://localhost:8080` to see location and a list of recommended POIs.
