package repository

// scannable is used to call .Scan for both sql.Row and sql.Rows
type scannable interface {
	Scan(dest ...interface{}) error
}
