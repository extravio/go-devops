# Golang for Devops

## Set up


```
# first run
docker run -v $(pwd):/usr/src/app --name godevops -it golang:1.20.1-bullseye /bin/bash

# then
docker container start godevops
docker exec -it godevops /bin/bash
```