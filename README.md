[![codecov](https://codecov.io/gh/DAtek/attiladudas/branch/main/graph/badge.svg?token=DD4YIGPFYE)](https://codecov.io/gh/DAtek/attiladudas)

# Source code of [attiladudas.com](https://attiladudas.com)

## Requirements:
- `just`
- `docker`
- `direnv`
- `go`
- `node`

## Running the tests locally:
- `just run-integration-test-services -d`
- `just migrate up`
- `just test`

## Running both the frontend and backend locally
- `just run-all-services -d`

On linux if you have problems with uploading files, then it's probably because of some docker permission errors. In that case recreate the `tmp/media` directory and permit everything (`chmod 777`), then restart everything and it'll work.
