# Source Code Annotation application

In order to evaluate quality of ML models, as well as to create “ImageNet for source core” there is a need for tools to automate the data collection/labeling/annotation.

## Installation

### Docker

```bash
docker build -t code-annotation .
docker run --rm -p 8080:8080 code-annotation
```

### Non-docker

```bash
go get <here will be path>
cd $GOPATH/<here will be path>
make serve
```

## Contributing

Please take a look at [CONTRIBUTING](CONTRIBUTING.md) file to see how to contribute in this project, get more information about the dashboard [architecture](CONTRIBUTING.md#Architecture) and how to launch it for [development](CONTRIBUTING.md#Development) purposes.

## License

GPLv3, see [LICENSE](LICENSE)
