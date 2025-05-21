# Docker Deployment

This directory contains Docker configuration files for running npm-download-spammer.

## Quick Start

### Build and run using Docker Compose

```bash
cd deploy
docker-compose up --build
```

### Run with predefined environment variables

Edit the `docker-compose.yml` file to uncomment and set your environment variables:

```yaml
environment:
  NPM_PACKAGE_NAME: "your-package-name"
  NPM_NUM_DOWNLOADS: 1000
  NPM_MAX_CONCURRENT_DOWNLOAD: 300
  NPM_DOWNLOAD_TIMEOUT: 3000
```

Then run:

```bash
docker-compose up --build
```

### Using Docker directly

Build the image:

```bash
docker build -t npm-download-spammer -f deploy/Dockerfile .
```

Run with interactive input:

```bash
docker run -it npm-download-spammer
```

Run with environment variables:

```bash
docker run -e NPM_PACKAGE_NAME="your-package-name" -e NPM_NUM_DOWNLOADS=1000 npm-download-spammer
```

## Configuration

The following environment variables are available:

| Variable | Description | Default |
|----------|-------------|---------|
| NPM_PACKAGE_NAME | NPM package name | (required) |
| NPM_NUM_DOWNLOADS | Number of downloads to add | 1000 |
| NPM_MAX_CONCURRENT_DOWNLOAD | Number of concurrent downloads | 300 |
| NPM_DOWNLOAD_TIMEOUT | Download timeout in milliseconds | 3000 | 