# First time setup
- Copy `.env.example` to `.env` and `.env.development` and fill out with the appropriate values based on your desired environment
- `docker-compose up`

# Preparing a production build
- `go build -ldflags="-X 'main.environment=production'"`
    - This replaces the `environment` variable within the `main` package and allows us to swap which env file to use based on the build type

# To run the final production executable
- `./go-web-service.exe`