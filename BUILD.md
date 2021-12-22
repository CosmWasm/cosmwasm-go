# Building docker images

Let's see if this works...

```
docker build -t demo/builder:latest -f Dockerfile .

# failure cases
docker run -v "$(pwd):/code" demo/builder:latest .
docker run -v "$(pwd):/code" demo/builder:latest foo

# RUN
docker run -v "$(pwd):/code" demo/builder:latest ./example/queue
docker run -e PAGES=30 -v "$(pwd):/code" demo/builder:latest ./example/queue

# FULL RUN
docker run -e CHECK=1 -e STRIP=1 -v "$(pwd):/code" demo/builder:latest ./example/queue


# DEBUG
docker run -it --entrypoint /bin/bash demo/builder:latest
```

