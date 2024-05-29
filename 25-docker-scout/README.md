## packagemain #25: Identifying Container Image vulnerabilities with Docker Scout

We all know, that Docker technology is great and brings us many advantages, but also, unfortunately, Docker images include many attack surfaces on different layers.

Every day, there are new vulnerabilities discovered in open source projects and maintainers are tasked with patching their software. [~30k new vulnerabilities discovered in 2023 alone](https://www.cvedetails.com/).

So how can we mitigate this risk? One solution is vulnerability scanning and its integration into your development lifecycle.

There are many free, open-source or paid tools for Docker vulnerability scanning:

- Docker Scout
- Aqua Security Trivy
- Snyk
- tenable.io
- and many others

They all have their advantages and differences, but still share the main goal: identify unpatched vulnerabilities.

Let's take a deep dive into the Docker vulnerability scanning and see it in action!

- Craft a sample Dockerfile as a foundation for our exploration.
- Scan for vulnerabilities with Docker Scout.
- Explore some resolution options.
- Set up a simple CI/CD pipeline to automate continuous scanning/reporting.

Our Dockerfile will use `golang:1.19` as a base image, which is not the latest version, but not so old either, and I believe many projects still use it.

```Dockerfile
FROM golang:1.19

WORKDIR /
COPY main.go .
RUN go build -o goapp main.go

CMD ["./goapp"]
```

In order to scan our image, we have to build it first, let's also run it to make sure it works as intended.

```bash
docker build -t goapp .
docker run goapp

> Hello, Docker!
```

If you're using a recent version of Docker Desktop, you might already have the Docker Scout command line tool available. While this guide won't endorse any specific tool, Docker Scout can be a handy starting point for this demonstration.

Some notes on Docker Scout:

- You need to have a Docker Hub account to run it
- [Free version has limitations](https://www.docker.com/products/docker-scout/) for remote images
- Not many registries are available, for example Google Cloud Container Registry is not available yet

```bash
docker login
docker scout enroll ORG_NAME
docker scout help
```

There are a few commands available:

- `quickview`: get a quick overview of an image, base image and available recommendations
- `compare`: compare an image to a second one (for instance to latest)
- `cves`: display vulnerabilities of an image
- `recommendations`: display available base image updates and remediation recommendations

We can use them to scan our already built image. Which resulted in...

```bash
docker scout cves goapp

80 vulnerabilities found in 25 packages
UNSPECIFIED  11
LOW          46
MEDIUM       12
HIGH         9
CRITICAL     2
```

Quite a lot right? Docker Scout can give us some basic recommendations with `recommendations` command. In our case it suggests to bump our base image to `golang:1.22`. Let's try it again.

```Dockerfile
FROM golang:1.22

WORKDIR /
COPY main.go .
RUN go build -o goapp main.go

CMD ["./goapp"]
```

```bash
docker build -t goapp .
docker scout cves goapp

59 vulnerabilities found in 23 packages
LOW       55
MEDIUM    1
HIGH      2
CRITICAL  1
```

Much better, but still...

Let's open this Critical [CVE-2024-32002](https://nvd.nist.gov/vuln/detail/CVE-2024-32002). As it turns out it's coming from Git software, but do we actually need Git in our image to run our application? To build probably yes, because we need to download modules and repositories, but since Go is compiled, we don't need Git after our program is compiled.

Multi-stage builds can help us here to separate build and run stages, which also has other benefits such as smaller image.

```Dockerfile
FROM golang:1.22 as builder

WORKDIR /
COPY main.go .
RUN go build -o goapp main.go

FROM alpine:latest
COPY --from=builder /goapp .

CMD ["./goapp"]
```

```bash
docker scout cves goapp

No vulnerable packages detected
```

Awesome, no vulnerabilities were identified! Or at least we may think so and have some peace of mind :)

While manual scanning is valuable, integrating vulnerability checks into your CI/CD pipeline is crucial for serious vulnerability management. This ensures automated scanning for every build, preventing malicious software from reaching production.

Docker Scout has a [GitHub Action](https://github.com/docker/scout-action) to run the Docker Scout CLI as part of your workflows.

Here is an example workflow (`.github/workflows/docker-scout.yaml`) which runs Docker Scout on every push and reports only Critical and High vulnerabilities as a comment. This actions requires authentication to Docker Hub, so we should add `DOCKERHUB_USERNAME` and `DOCKERHUB_TOKEN` to secrets.

```yaml
name: Docker Scout

on:
  push:
    branches:
      - "*"

jobs:
  build-and-scan:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout repository
        uses: actions/checkout@v4
      - name: Setup Docker buildx
        uses: docker/setup-buildx-action@v3
      - name: Build Docker image
        uses: docker/build-push-action@v4.0.0
        with:
          context: .
          push: false
          load: true
          tags: ${{ github.event.repository.name }}
      - name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      - name: Docker Scout
        id: docker-scout
        uses: docker/scout-action@v1
        with:
          command: cves
          image: ${{ github.event.repository.name }}
          ignore-unchanged: true
          only-severities: critical,high
          write-comment: true
          github-token: ${{ secrets.GITHUB_TOKEN }}
```

That's it for today.

Incorporating these practices into your workflow empowers developers to streamline vulnerability management and maintain a more secure containerized ecosystem.
