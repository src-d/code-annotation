# Source Code Annotation application

In order to evaluate quality of ML models, as well as to create “ImageNet for source core” there is a need for tools to automate the data collection/labeling/annotation.

## Installation

### Github OAuth tokens

First you need OAuth application on github. [Read how to create it](https://developer.github.com/apps/building-oauth-apps/creating-an-oauth-app/).

On a [page](https://github.com/settings/developers) with your application you will need `Client ID` and `Client Secret`

Copy `.env.tpl` to `.env` and set tokens there.

### Docker

```bash
docker build -t srcd/code-annotation .
docker run --env-file .env --rm -p 8080:8080 srcd/code-annotation
```

### Non-docker

```bash
go get github.com/src-d/code-annotation/...
cd $GOPATH/github.com/src-d/code-annotation
make serve
```

## Development

Frontend:

```
yarn
yarn start
```

## Contributing

Please take a look at [CONTRIBUTING](CONTRIBUTING.md) file to see how to contribute in this project, get more information about the dashboard [architecture](CONTRIBUTING.md#Architecture) and how to launch it for [development](CONTRIBUTING.md#Development) purposes.

## License

GPLv3, see [LICENSE](LICENSE)
