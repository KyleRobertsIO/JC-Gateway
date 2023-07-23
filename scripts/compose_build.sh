#!/bin/bash

{
    docker image rm jc-gateway-job-manager --force
} || {
    echo "No existing docker image build to remove"
}
docker compose up