[![Tests](https://github.com/DI-Tony-Reed/go-web-service/actions/workflows/tests.yml/badge.svg)](https://github.com/DI-Tony-Reed/go-web-service/actions/workflows/tests.yml)
![Coverage Badge](https://img.shields.io/badge/Coverage-91.6%25-brightgreen.svg)

# First time setup
- DDo the following within the `frontend` and `server` directories
  - Copy `.env.example` to `.env.production` and `.env.development` and fill out with the appropriate values based on your
    desired environment
- `docker-compose up --build`
    - The MySQL container will be built and ran first
    - The backend container uses Air to watch for files and automatically recompile the application; this container
      will not attempt to run until the MySQL image is running and returns a healthy response
    - The frontend container includes Vite, Vue, and other dependencies and will watch for changes as well
- Visit `localhost:8000` to interact with the front end

# Building backend executables

- `make build`

# To run generated executables

- `./server/bin/go-web-service-{ENVIRONMENT}`
