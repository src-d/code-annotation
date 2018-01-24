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

## Importing and Exporting Data

### Import File Pairs for Annotation

The file pairs must be initially provided via a database, either [SQLite](https://sqlite.org/) or [PostgreSQL](https://www.postgresql.org/). The database should contain rows of file pairs in the following table:

```SQL
CREATE TABLE files (
    blob_id_a TEXT, repository_id_a TEXT, commit_hash_a TEXT, path_a TEXT, content_a TEXT,
    blob_id_b TEXT, repository_id_b TEXT, commit_hash_b TEXT, path_b TEXT, content_b TEXT,
    score DOUBLE PRECISION);
```

The `import` command will use those file pairs to create a new database that will be used internally by the Annotation Tool. The destination database does not need to be empty, new imported file pairs can
be added to previous imports.

_Please note_: if a file pair is identical to an existing one it will not be detected. A new pair entry will be created with the same contents.

To use it, run it as:

```bash
$ import <origin-DSN> <destination-DSN>
```

Where both `DSN` (Data Source Name) arguments must be one of:

* `sqlite:///path/to/db.db`
* `postgresql://[user[:password]@][netloc][:port][,...][/dbname]`

Some usage examples:
```bash
$ import sqlite://./input.db sqlite:///home/user/input.db

$ import sqlite://input.db postgres://testing:testing@localhost:5432/input?sslmode=disable
```

For a complete reference of the PostgreSQL connection string, see the [documentation for the lib/pq Go package](https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters).

### Export Annotation Results

To work with the annotation results, the internal data can be extracted in a new database using the `export` command.

```bash
$ export <origin-DSN> <destination-DSN>
```

The arguments use the same format as the `import` tool, see the previous section.

In this case, origin will be the internal database, and destination the new database. The new database will have the same contents are the internal database.

To study the user annotation results, focus on the **`assignments`** table.


## Contributing

[Contributions](https://github.com/src-d/code-annotation/issues) are more than welcome, if you are interested please take a look to
our [Contributing Guidelines](CONTRIBUTING.md).

# Code of Conduct

All activities source{d} projects are governed by the [source{d} code of conduct](CODE_OF_CONDUCT.md).

## License

GPLv3, see [LICENSE](LICENSE)
