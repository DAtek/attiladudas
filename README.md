[![codecov](https://codecov.io/gh/DAtek/attiladudas/graph/badge.svg?token=DD4YIGPFYE)](https://codecov.io/gh/DAtek/attiladudas)

# Source code of [attiladudas.com](https://attiladudas.com)

## Requirements:
- `just`
- `docker`
- `direnv`
- `go`
- `gotestsum`
- `node`
- `awscli`
- For the backend deployment: `ansible`

## Running the tests locally:
- `cd backend`
- `just test`

## Running both the frontend and backend locally
- `docker compose up --build`
