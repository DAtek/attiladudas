coverfile := "tmp/.coverage"
pkgs := "\
./backend/app/gallery \
./backend/app/auth \
./backend/app/five_in_a_row \
./backend/components \
./backend/components/auth \
./backend/components/gallery \
./backend/components/five_in_a_row \
./backend/components/room_manager "


test *options:
    go test {{ options }} {{ pkgs }}


test-cover *options:
    go test {{ options }} -coverprofile {{ coverfile }} {{ pkgs }}


show-coverage:
    go tool cover -html={{ coverfile }}


test-and-show-covarage: test-cover show-coverage


run-all-services *args:
    docker-compose --profile all up --build {{ args }}


reload-backend:
    docker-compose up --build -d backend


stop-all-services:
    docker-compose --profile all stop


run-integration-test-services *args:
    docker-compose --profile integration_test up {{ args }}


stop-integration-test-services:
    docker-compose --profile integration_test stop


generate-keypair:
    openssl genpkey -algorithm ed25519 -out tmp/private.pem
    openssl pkey -in tmp/private.pem -pubout -out tmp/public.pem


migrate *args:
    go run ./backend/cmd/migrate {{ args }}


create-user *args: 
    go run ./backend/cmd/create_user {{ args }}


build: clean build-backend build-frontend

clean:
    rm -rf .dist

build-backend:
    #!/bin/bash

    echo "Building backend"
    mkdir .dist 2>/dev/null
    rm -rf .dist/backend
    mkdir .dist/backend
    for item in create_user migrate server
    do
        go build ./backend/cmd/${item}
        mv ${item} .dist/backend/${item}
    done
    cp -r backend/migrations .dist/backend/.

build-frontend:
    #!/bin/bash

    echo "Building frontend"
    mkdir .dist 2>/dev/null
    cd frontend && npm run build
