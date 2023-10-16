# First time setup
- Copy `.env.example` to `.env.production` and `.env.development` and fill out with the appropriate values based on your desired environment
- `docker-compose up`

# Preparing a production build
- (_If Go is installed locally_) `go build -ldflags="-X 'go-web-service/src/utils.environment=production'"`
- Otherwise, create the executable from our container: `docker exec app go build -ldflags="-X 'go-web-service/src/utils.environment=production'"`
  - _A note on this, if you are on a Windows machine, this generated binary from the container will not be executable locally_
- The `-ldflags="-X...` command replaces the `environment` variable within the specified package and allows us to swap which env file to use based on the build type

# To run the final production executable
- `./go-web-service`