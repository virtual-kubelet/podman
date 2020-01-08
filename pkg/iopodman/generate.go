package iopodman

// For this generator to work you need remote podman setup
// https://podman.io/blogs/2019/01/16/podman-varlink.html

//go:generate go get -u github.com/varlink/go/cmd/varlink-go-interface-generator
//go:generate go install github.com/varlink/go/cmd/varlink-go-interface-generator
//go:generate $GOPATH/bin/varlink-go-interface-generator io.podman.varlink
