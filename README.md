# First time setup
- Copy `.env.example` to `.env.production` and `.env.development` and fill out with the appropriate values based on your desired environment
- `docker-compose up`

# Preparing a production build
- `go build -ldflags="-X 'go-web-service/src/utils.environment=production'"`
    - This replaces the `environment` variable within the specified package and allows us to swap which env file to use based on the build type

# To run the final production executable
- `./go-web-service.exe`