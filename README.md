[![Tests](https://github.com/DI-Tony-Reed/go-web-service/actions/workflows/tests.yml/badge.svg)](https://github.com/DI-Tony-Reed/go-web-service/actions/workflows/tests.yml)
![Coverage Badge](https://img.shields.io/badge/Coverage-91.6%25-brightgreen.svg)

# First time setup

- Copy `.env.example` to `.env.production` and `.env.development` and fill out with the appropriate values based on your
  desired environment
- Do the same for the `.env.example` file within the `frontend` directory
- `docker-compose up --build`
    - The MySQL container will be built and ran first
    - The backend-development container uses Air to watch for files and automatically recompile the application; this container
      will not attempt to run until the MySQL image is running and returns a healthy response
    - The frontend-development container includes Vite, Vue, and other dependencies and will watch for changes
- Visit `localhost:8000` to interact with the front end

# Preparing a production build

- `docker exec -it backend-development bash`
- `make build`
- `exit`

# To run the final production executable

- `./go-web-service-{ENVIRONMENT}`
