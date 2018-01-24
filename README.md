[![Build Status](https://travis-ci.org/src-d/go-git.svg)](https://travis-ci.org/src-d/go-git)
![unstable](https://svg-badge.appspot.com/badge/stability/unstable?a)

# Source Code Annotation Tool

In order to evaluate quality of ML models, as well as to create “ImageNet for source core” there is a need for tools to automate the data collection/labeling/annotation.

![Screenshot](.github/screenshot.png?raw=true)

## Installation

### Github OAuth tokens

1. You need OAuth application on github. [Read how to create it](https://developer.github.com/apps/building-oauth-apps/creating-an-oauth-app/).

    `Authorization callback URL: http://127.0.0.1:8080/oauth-callback`

2. Copy `.env.tpl` to `.env`.

3. On a [page](https://github.com/settings/developers) with your application find `Client ID` and `Client Secret` and put them in `.env` file.

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

## Import files for annotation

TODO

## Contributing

[Contributions](https://github.com/src-d/code-annotation/issues) are more than welcome, if you are interested please take a look to
our [Contributing Guidelines](CONTRIBUTING.md).

# Code of Conduct

All activities source{d} projects are governed by the [source{d} code of conduct](CODE_OF_CONDUCT.md).

## License

GPLv3, see [LICENSE](LICENSE)
