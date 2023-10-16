# First time setup
- Copy `.env.example` to `.env.production` and `.env.development` and fill out with the appropriate values based on your desired environment
- `docker-compose up`
  - The MySQL container will be built and ran first
  - The development container uses Air to watch for files and automatically recompile the application; this container will not attempt to run until the MySQL image is running and returns a healthy response

# Preparing a production build
- (_If Go is installed locally_) `go build -ldflags="-X 'main.environment=production'"`
- Otherwise, create the executable from our container: `docker exec development go build -ldflags="-X 'main.environment=production'"`
  - _A note on this, if you are on a Windows machine, this generated binary from the container will not be executable locally_
- The `-ldflags="-X...` command replaces the `environment` variable within the specified package and allows us to swap which env file to use based on the build type

# To run the final production executable
- `./go-web-service`

TODO TEMP
docker build -t production - < Docker/Dockerfile.production 
