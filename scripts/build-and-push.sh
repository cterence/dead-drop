#!/bin/bash
docker build -t ghcr.io/cterence/dead-drop:latest .
docker push ghcr.io/cterence/dead-drop:latest
