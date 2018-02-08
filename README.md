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

## Installation

## Github OAuth tokens

1. You need an OAuth application on GitHub. See [how to create OAuth applications on GitHub](https://developer.github.com/apps/building-oauth-apps/creating-an-oauth-app/).

In order to be able to use this application while running the tool locally, make sure you add http://127.0.0.1:8080/oauth-callback to the authorization callback URL field.

2. Copy `.env.tpl` to `.env`.

3. Retrieve the values for your application's Client ID and Client Secret from the [GitHub Developer Settings page](https://github.com/settings/developers) and add them to the end of the corresponding lines in .env.

### Docker

docker run --env-file .env --rm -p 8080:8080 srcd/code-annotation

### Non-docker

Download binary from [releases](https://github.com/src-d/code-annotation/releases) for your platform.

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

#### Set the Internal Database Connection

Before starting the application you will need to set the `DB_CONNECTION` variable in the `.env` file. It should point to the database created with the `import` command.

This variable uses the same `DSN` string as the `import` command to point to a SQLite or PosgreSQL database.

Some examples:

```
DB_CONNECTION=sqlite:///home/user/internal.db
```

```
DB_CONNECTION=postgres://testing:testing@localhost:5432/input?sslmode=disable
```

### Export Annotation Results

To work with the annotation results, the internal data can be extracted into a new SQLite database using the `export` command.

```bash
$ export <origin-DSN> <path-to-sqlite.db>
```

The DSN argument uses the same format as the `import` tool, see the previous section.

In this case, origin will be the internal database, and destination the new database. This new database will have the same contents as the internal one.

The annotations made by the users will be stored in the **`assignments`** table.

## Access control

It is possible to restrict access and choose each user's role by adding their GitHub accounts to a specific [organization](https://help.github.com/articles/collaborating-with-groups-in-organizations/) or [team](https://help.github.com/articles/organizing-members-into-teams/).

This is optional, if you don't set any restrictions all users with a valid GitHub account will be able to login as a Requester. You may also set a restriction only for Requester users, and leave open access to anyone as Workers.

To do so, set the following variables in your `.env` file:

* `OAUTH_RESTRICT_ACCESS`
* `OAUTH_RESTRICT_REQUESTER_ACCESS`

Both variables accept a string with either `org:<organization-name>` or `team:<team-id>`. For example:

```bash
OAUTH_RESTRICT_ACCESS=org:my-organization
OAUTH_RESTRICT_REQUESTER_ACCESS=team:123456
```

## Contributing

[Contributions](https://github.com/src-d/code-annotation/issues) are more than welcome, if you are interested please take a look to our [Contributing Guidelines](CONTRIBUTING.md).

# Code of Conduct

All activities under source{d} projects are governed by the [source{d} code of conduct](CODE_OF_CONDUCT.md).

## License

GPLv3, see [LICENSE](LICENSE)
