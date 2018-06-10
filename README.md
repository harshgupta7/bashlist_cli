# Bashlist CLI
CLI Package for Bashlist.


# Setup
* Install go
* Install external dependencies using `go get <dependency-name>`.
* Build Using `go build cli-mac.go`--This will create an executable binary.
* Run `./cli-mac`

# docs/
* Contains CLI docs for customers
* run `npm i docsify-cli -g` to install docsify
* run `docsify serve ./docs` -- This will start a local server
* open `localhost:3000` to see full CLI documentation.

## Dependencies

* "github.com/buger/jsonparser"
* "github.com/docker/docker-credential-helpers/credentials"
* "github.com/docker/docker-credential-helpers/osxkeychain"
* "github.com/fatih/color"
* "github.com/howeyc/gopass"
* "github.com/imroc/req"
* "github.com/olekukonko/tablewriter"
* "github.com/skratchdot/open-golang/open"
