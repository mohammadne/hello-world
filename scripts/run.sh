#!/bin/bash

hello_world_version="v0.1.1"
output_directory="./outputs"

podman run -d --name my_container ghcr.io/mohammadne/hello-world:$hello_world_version
podman exec -it my_container /app/entrypoint executer
podman cp my_container:/app/outputs $output_directory
podman rm -f my_container -t 0
