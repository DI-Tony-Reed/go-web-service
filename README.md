[![Tests](https://github.com/DI-Tony-Reed/go-web-service/actions/workflows/tests.yml/badge.svg)](https://github.com/DI-Tony-Reed/go-web-service/actions/workflows/tests.yml)
![Coverage Badge](https://img.shields.io/badge/Coverage-93.7%25-brightgreen.svg)

# First time setup
- Copy `.env.example` to `.env` and fill out with the appropriate values based on your
  desired environment within the `server` and `frontend` dirs
- `docker-compose up --build`
    - The MySQL container will be built and ran first
    - The backend container uses Air to watch for files and automatically recompile the application; this container
      will not attempt to run until the MySQL image is running and returns a healthy response
    - The frontend container includes Vite, Vue, and other dependencies and will watch for changes as well
- Visit `localhost:8000` to interact with the front end
