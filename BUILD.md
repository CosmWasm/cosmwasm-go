# Building docker images

Let's see if this works...

```
docker build -t demo/builder:latest -f Dockerfile.arm .

docker run -it demo/builder:latest /bin/bash
```

