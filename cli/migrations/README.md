### Disclaimer:

These scripts are <u>**not real migrations**</u>. There is no up/down rules. These scripts are not production ready. The idempotence of these scripts is not guaranteed at all. The internal DB should be SQLite.

These scripts were needed to migrate the Database between different states.

- [Vacuum database](command-vacuum.go)
```shell
# rebuilds the database to defragment it
go run cli/migrations/*.go vacuum internal.db
```

- [Add UAST nullable columns](command-UAST-add-columns.go)
```shell
# prepares the current "internal.db". Adds "uast_a" and "uast_b" nulable cols
go run cli/migrations/*.go uast-add-cols internal.db
```

- [Add UAST to a database, reading from other Database](command-UAST-import-from-source-db.go)
```shell
# import UASTs into "internal.db" reading from "source.db"
go run cli/migrations/*.go uast-import internal.db source.db
```

- [Remove diff column from a database](command-diff-remove-column.go)
```shell
# prepares the current "internal.db". Remove "diff" col
go run cli/migrations/*.go diff-rm-col internal.db
```
