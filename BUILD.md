# Building docker images

Let's see if this works...

```
VERSION=0.3.0-arm64
docker build -t cosmwasm/go-optimizer:${VERSION} -f Dockerfile .

# failure cases
docker run -v "$(pwd):/code" cosmwasm/go-optimizer:${VERSION} .
docker run -v "$(pwd):/code" cosmwasm/go-optimizer:${VERSION} foo

# RUN
docker run -v "$(pwd):/code" cosmwasm/go-optimizer:${VERSION} ./example/queue
docker run -e PAGES=30 -v "$(pwd):/code" cosmwasm/go-optimizer:${VERSION} ./example/queue

# FULL RUN
docker run -e CHECK=1 -e STRIP=1 -v "$(pwd):/code" cosmwasm/go-optimizer:${VERSION} ./example/queue


# DEBUG
docker run -it --entrypoint /bin/bash cosmwasm/go-optimizer:${VERSION}
```

