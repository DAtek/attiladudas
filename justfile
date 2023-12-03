run-all *args:
    docker-compose up --build {{ args }}


watch-backend:
    #!/bin/bash

    inotifywait -q -m -e close_write -r ./backend |
    while read input; do
        docker-compose stop backend
        docker-compose up --build backend &
    done


stop-all-services:
    docker-compose stop


migrate *args:
    just -f ./backend/justfile migrate {{ args }}


create-user *args: 
    just -f ./backend/justfile create-user {{ args }}


build: clean build-backend build-frontend


clean:
    rm -rf .dist


build-backend:
    #!/bin/bash

    echo "Building backend"
    mkdir .dist 2>/dev/null
    rm -rf .dist/backend
    mkdir .dist/backend
    export CGO_ENABLED=0
    for item in create-user migrate run-server
    do
        go build -C backend ./cmd/${item}
        mv backend/${item} .dist/backend/${item}
    done
    cp -r backend/db/migrations .dist/backend/.
    echo "Done"

build-frontend:
    #!/bin/bash

    echo "Building frontend"
    mkdir .dist 2>/dev/null
    rm -rf .dist/frontend
    cd frontend && pnpm run lint && pnpm run build
