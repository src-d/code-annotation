[![Build Status](https://travis-ci.org/src-d/code-annotation.svg)](https://travis-ci.org/src-d/code-annotation)
![unstable](https://svg-badge.appspot.com/badge/stability/unstable?a)

# Source Code Annotation Tool

Training Machine Learning models often requires large datasets to be duly annotated.
The nature of these annotations vary depending on the dataset considered: they can be
the number to be recognized in the [MNIST dataset](http://yann.lecun.com/exdb/mnist/),
the coordinates of the box containing the objects to be identified in an object detection problem, etc.

This tool provides a simple UI to add annotations to existing datasets, a command line tool
to fetch more elements to be annotated, and an export mechanism.

Currently, the project provides one single example consisting on labeling two pieces of code
as being identical, similar, or different.

Source code annotation tool offers an UI to annotate source code and review these annotations, and a CLI to define the code to be annotated and export the annotations.

![Screenshot](.github/screenshot.png?raw=true)

## Requirements

### Global dependencies

You should already have [Go installed](https://golang.org/doc/install#install), and properly [configured the $GOPATH](https://github.com/golang/go/wiki/SettingGOPATH)

```
go version; # prints your go version
echo $GOPATH; # prints your $GOPATH path
```

The project must be under $GOPATH, as required by the Go tooling.
You should be able to navigate into the source code by running:

```
cd $GOPATH/src/github.com/src-d/code-annotation
```

You need also [Yarn v1.x.x installed](https://yarnpkg.com/en/docs/install)

```
yarn --version; # prints your Yarn version
```

### Github OAuth tokens

1. You need an OAuth application on GitHub. See [how to create OAuth applications on GitHub](https://developer.github.com/apps/building-oauth-apps/creating-an-oauth-app/).

In order to be able to use this application while running the tool locally, make sure you add http://127.0.0.1:8080/oauth-callback to the authorization callback URL field.

2. Copy `.env.tpl` to `.env`.

3. Retrieve the values for your application's Client ID and Client Secret from the [GitHub Developer Settings page](https://github.com/settings/developers) and add them to the end of the corresponding lines in .env.

## Installation

You need to satisfy all [project requirements](#requirements), and then run:

```bash
$ go get github.com/src-d/code-annotation/...
$ cd $GOPATH/github.com/src-d/code-annotation
$ make serve
```

This will start a server locally, which you can access on [http://localhost:8080](http://localhost:8080)

## Importing and Exporting Data

### Import File Pairs for Annotation

The file pairs must be provided via an [SQLite](https://sqlite.org/) database. The database **must follow the expected schema**, please [follow this link](./cli/examples/import/example.sql) to see an example.

The `import` command will use those file pairs to create a new [SQLite](https://sqlite.org/) or [PostgreSQL](https://www.postgresql.org/) database that will be used internally by the Annotation Tool. The destination database does not need to be empty, new imported file pairs can be added to previous imports.

_Note_: duplicate entries are not filtered, so running an import multiple times will result in repeated rows.

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
Imported 989 file pairs successfully

$ import /home/user/input.db postgres://testing:testing@localhost:5432/input?sslmode=disable
Imported 562 file pairs successfully
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

All activities under source{d} projects are governed by the [source{d} code of conduct](CODE_OF_CONDUCT.md).

## License

GPLv3, see [LICENSE](LICENSE)
