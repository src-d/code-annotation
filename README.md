[![Build Status](https://travis-ci.org/src-d/code-annotation.svg)](https://travis-ci.org/src-d/code-annotation)
![unstable](https://svg-badge.appspot.com/badge/stability/unstable?a)

# Source Code Annotation Tool

Source code annotation tool offers an UI to annotate source code and review these annotations, and a CLI to define the code to be annotated and export the annotations.

### Why is is needed?

In order to evaluate quality of ML models, as well as to create “ImageNet for Source Code” there is a need for tools to automate the data collection/labeling/annotation.

![Screenshot](.github/screenshot.png?raw=true)

## Requirements

### Global dependencies

You should already have [Go installed](https://golang.org/doc/install#install), and properly [configured the $GOPATH](https://github.com/golang/go/wiki/SettingGOPATH)

```
go version; # prints your go version
echo $GOPATH; # prints your $GOPATH path
```

The project must be under the `$GOPATH`, following the Go import conventions, what means you can go to its directory running:

```
cd $GOPATH/src/github.com/src-d/landing
```

You need also [Yarn v1.x.x installed](https://yarnpkg.com/en/docs/install)

```
yarn --version; # prints your Yarn version
```

### Github OAuth tokens

1. You need an OAuth application on GitHub. See [how to create OAuth applications on GitHub](https://developer.github.com/apps/building-oauth-apps/creating-an-oauth-app/).

    `Authorization callback URL: http://127.0.0.1:8080/oauth-callback`

2. Copy `.env.tpl` to `.env`.

3. Retrieve the values for your application's Client ID and Client Secret from the [GitHub Developer Settings page](https://github.com/settings/developers) and add them to the end of the corresponding lines in .env.

## Installation

You need to satisfy all [project requirements](#requirements), and then run:

```bash
$ go get github.com/src-d/code-annotation/...
$ cd $GOPATH/github.com/src-d/code-annotation
$ make serve
```

It runs everything you need to get the tool working at [http://localhost:8080](http://localhost:8080)

## Importing and Exporting Data

### Import File Pairs for Annotation

The file pairs must be initially provided via an [SQLite](https://sqlite.org/) database. The database **must follow the expected schema**, please [follow this link](./cli/examples/import/example.sql) to see an example.

The `import` command will use those file pairs to create a new [SQLite](https://sqlite.org/) or [PostgreSQL](https://www.postgresql.org/) database that will be used internally by the Annotation Tool. The destination database does not need to be empty, new imported file pairs can be added to previous imports.

_Please note_: if a file pair is identical to an existing one it will not be detected. A new pair entry will be created with the same contents.

To use it, run it as:

```bash
$ import <path-to-sqlite.db> <destination-DSN>
```

Where the `DSN` (Data Source Name) argument must be one of:

* `sqlite:///path/to/db.db`
* `postgresql://[user[:password]@][netloc][:port][,...][/dbname]`

Some usage examples:

```bash
$ import ./input.db sqlite:///home/user/internal.db

$ import /home/user/input.db postgres://testing:testing@localhost:5432/input?sslmode=disable
```

For a complete reference of the PostgreSQL connection string, see the [documentation for the lib/pq Go package](https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters).

### Export Annotation Results

To work with the annotation results, the internal data can be extracted into a new SQLite database using the `export` command.

```bash
$ export <origin-DSN> <path-to-sqlite.db>
```

The DSN argument uses the same format as the `import` tool, see the previous section.

In this case, origin will be the internal database, and destination the new database. This new database will have the same contents as the internal one.

To study the user annotation results, focus on the **`assignments`** table.

## Contributing

[Contributions](https://github.com/src-d/code-annotation/issues) are more than welcome, if you are interested please take a look to
our [Contributing Guidelines](CONTRIBUTING.md).

# Code of Conduct

All activities source{d} projects are governed by the [source{d} code of conduct](CODE_OF_CONDUCT.md).

## License

GPLv3, see [LICENSE](LICENSE)
