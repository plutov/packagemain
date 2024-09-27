## packagemain #26: Self-Hosting a Container Registry

## What is a Container Image?

Before we talk about container registries, let's first understand what a container image is. In a nutshell a container image is a package that includes all of the files, libraries, and configurations to run a container. They are composed of [layers](https://docs.docker.com/get-started/docker-concepts/building-images/understanding-image-layers/). Each layer represents a set of file system changes that add, remove, or modify files.

The most common way to create a container image is to use a `Dockerfile`.

```bash
docker build -t pliutau/hello-world:v0 .

# check the images locally
docker images
# REPOSITORY    TAG       IMAGE ID       CREATED          SIZE
# hello-world   latest    9facd12bbcdd   22 seconds ago   11MB
```

This creates a container image that is stored on your local machine. But what if you want to share this image with others or use it on a different machine? This is where container registries come in.

## What is a Container Registry?

A container registry is a storage catalog where you can push and pull container images from. The images are grouped into repositories, which are collections of related images with the same name. For example, on Docker Hub registry, [nginx](https://hub.docker.com/_/nginx) is the name of the repository that contains different versions of the nginx images.

Some registries are public, meaning that the images hosted on them are accessible to anyone on the Internet. Public registries such as [Docker Hub](https://hub.docker.com) are good option to host open-source projects.

On the other hand private registries provide a way to incorporate security and privacy into enterprise container image storage, either hosted in cloud or on-premises. These private registries often come with advanced security features and technical support. There is a growing list of private registries available such as [Amazon ECR](https://aws.amazon.com/ecr/), [GCP Artifact Registry](https://cloud.google.com/artifact-registry/docs), [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry), also Docker Hub offers a private repository feature.

As developer you interact with a container registry when using the `docker push` and `docker pull` commands.

```bash
docker push docker.io/pliutau/hello-world:v0

# In case of Docker Hub we could also skip the registry part
docker push pliutau/hello-world:v0
```

Let's look at the anatomy of a container image URL.

```bash
docker pull docker.io/pliutau/hello-world:v0@sha256:dc11b2...
                |            |            |          |
                ↓            ↓            ↓          ↓
             registry    repository      tag       digest
```

## Why you might want to self-host a Container Registry?

Sometimes, instead of relying on a provider, like AWS or GCP, you might want to host your images yourself. This keeps your infrastructure internal and makes you less reliant on external vendors. In some heavily regulated industries, this is even a requirement.

A self-hosted registry runs on your own servers, giving you more control over how the registry is configured and where the container images are hosted.

## How to self-host a Container Registry?

There are several open-source container registry solutions available. The most popular one is officially supported by Docker, called [registry](https://hub.docker.com/_/registry) with its implementation for storing and distributing of container images and artifacts. Which means that you can run your own registry inside a container.

Here are the main steps to run a registry on a server:
- Install Docker and Docker Compose on the server
- Configure and run the **registry** container
- Run **Nginx** for handling TLS and forwarding requests to the registry container
- Setup certificates and configure a domain

## Server

You can use any server that supports Docker. For example, you can use a DigitalOcean Droplet with Ubuntu 20.04. For this demo I assume that you already have a server. I used Google Cloud Compute to create a VM with Ubuntu.

Once we're inside our VM we should install Docker and Docker Compose. Docker Compose is optional, but it makes it easier to manage multi-container applications.

```bash

```bash
# install docker engine
sudo apt install docker.io
sudo docker --version

# install docker-compose
# look here for the latest version and arch - https://github.com/docker/compose/releases
sudo curl -L "https://github.com/docker/compose/releases/download/v2.29.7/docker-compose-linux-aarch64" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
sudo docker-compose --version
```

## Registry Configuration

Next we need to configure our registry container. The following **compose.yaml** file will create a registry container with a volume for storing the images and a volume for storing the password file.

```yaml
services:
  registry:
    image: registry:latest
    environment:
      REGISTRY_AUTH: htpasswd
      REGISTRY_AUTH_HTPASSWD_REALM: Registry Realm
      REGISTRY_AUTH_HTPASSWD_PATH: /auth/registry.password
      REGISTRY_STORAGE_FILESYSTEM_ROOTDIRECTORY: /data
    volumes:
      # Mount the password file
      - ./registry/registry.password:/auth/registry.password
      # Mount the data directory
      - ./registry/data:/data
    ports:
      - 5000
```

The password file is used to authenticate users when they push or pull images from the registry. You can create a password file using the `htpasswd` command. We also should create a folder for storing the images.

```bash
cd ~
mkdir -p ./registry/data

# create a password file
mkdir registry
htpasswd -Bbn busy bee > ./registry/registry.password
```

Now we can start the registry container.

```bash
docker-compose up

# successfull run should output something like this:
# registry | level=info msg="listening on [::]:5000"
```

## Nginx

As mentioned earlier, we can use Nginx to handle TLS and forward requests to the registry container. We can add nginx service to our **compose.yaml** file.

```yaml
services:
  registry:
    # ...
  nginx:
    image: nginx:latest
    depends_on:
      - registry
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./nginx/certs:/etc/nginx/certs
    ports:
      - "443:443"
```

The **nginx.conf** file could look like this:

```ini

```

## Other options

Following the example above you can also run the registry on Kubernetes. Or you could use a managed registry service like [Harbor](https://goharbor.io/), which is an open-source registry that provides advanced security features and is compatible with Docker and Kubernetes.

Also, if you want to have a UI for your self-hosted registry you could use a project like [Joxit/docker-registry-ui](https://github.com/Joxit/docker-registry-ui) and run it in a separate container.

## Conclusion

Self-hosted Container Registries allow you to have complete control over your registry and the way it's deployed.

Whatever your reasons for running a self-hosted registry, you now know how it's done. From here you can compare the different options and choose the one that best fits your needs.
