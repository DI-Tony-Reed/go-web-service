# First time setup
- Copy `.env.example` to `.env.production` and `.env.development` and fill out with the appropriate values based on your desired environment
- `docker-compose up`
  - The MySQL container will be built and ran first
  - The development container uses Air to watch for files and automatically recompile the application; this container will not attempt to run until the MySQL image is running and returns a healthy response

# Preparing a production build
- `make build`

# To run the final production executable
- `./go-web-service-{ENVIRONMENT}`
