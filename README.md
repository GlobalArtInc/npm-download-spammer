# npm-download-spammer

A tool for increasing NPM package download counters.

## Installation

### From source code

```bash
git clone https://npm-download-spammer.git
cd npm-download-spammer
go build -o npm-download-spammer
```

## Usage

### Interactive mode

Run the program without arguments:

```bash
./npm-download-spammer
```

You will be prompted to enter:
- Package name
- Number of downloads
- Maximum number of concurrent downloads
- Download timeout (in milliseconds)

### Using environment variables

```bash
NPM_PACKAGE_NAME="your-package" \
NPM_NUM_DOWNLOADS=1000 \
NPM_MAX_CONCURRENT_DOWNLOAD=300 \
NPM_DOWNLOAD_TIMEOUT=3000 \
./npm-download-spammer
```

### Using configuration file

Create a `npm-downloads-increaser.json` file in the launch directory:

```json
{
    "packageName": "your-package",
    "numDownloads": 1000,
    "maxConcurrentDownloads": 300,
    "downloadTimeout": 3000
}
```

Then run the program:

```bash
./npm-download-spammer
```

## Configuration parameters

| Parameter | Description | Default |
|----------|----------|--------------|
| packageName | NPM package name | (required) |
| numDownloads | Number of downloads to add | 1000 |
| maxConcurrentDownloads | Number of concurrent downloads | 300 |
| downloadTimeout | Download timeout in milliseconds | 3000 |

## Notes

- For slow connections, it's recommended to reduce `maxConcurrentDownloads` and increase `downloadTimeout`
- The program works with scoped packages, for example `@scope/package-name`

## Docker

For Docker deployment instructions, see [deploy/README.md](deploy/README.md).

## License

MIT
