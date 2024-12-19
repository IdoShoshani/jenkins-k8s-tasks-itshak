#!/bin/bash
echo "DOCKER STARTS HERE"
service --status-all
service docker start
service docker start
docker version
docker ps
echo "DOCKER ENDS HERE"
sleep 100000
