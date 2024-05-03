#!/bin/bash

hello_world_version="v0.1.2"
output_directory="./outputs"

docker run -d --name hello_world ghcr.io/mohammadne/hello-world:$hello_world_version
docker exec -it hello_world /app/entrypoint executer
docker cp hello_world:/app/outputs $output_directory
docker rm -f hello_world -t 0

server_xray_directory="/root/xray"

ssh proxy-server -t 'rm -rf '"$server_xray_directory"''
scp -r "$output_directory/server" proxy-server:$server_xray_directory
ssh proxy-server -t 'cd '"$server_xray_directory"' && docker compose up -d'

cd "$output_directory/client" && docker compose up -d
