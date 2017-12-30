### How to build a tiny Docker image for Go program

```
docker build -t packagemain .
docker build -f Dockerfile.tiny -t packagemain-tiny .
docker images packagemain
docker images packagemain-tiny
docker run packagemain
docker run packagemain-tiny
```