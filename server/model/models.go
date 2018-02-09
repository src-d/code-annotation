package model

import (
	"database/sql"
	"database/sql/driver"
	"errors"
)

// User of the application; can be Requester or Workers
type User struct {
	ID        int
	Login     string // GitHub account username
	Username  string // Real name, as returned by GitHub
	AvatarURL string
	Role      Role
}

// Experiment groups a certain amount of FilePairs
type Experiment struct {
	ID          int
	Name        string
	Description string
}

// Assignment tracks the answer of a worker to a given FilePair of an Experiment
type Assignment struct {
	ID           int
	UserID       int
	PairID       int
	ExperimentID int
	Answer       sql.NullString
	Duration     int
}

// FilePair represents the pairs of files to annotate
type FilePair struct {
	ID           int
	Score        float64
	Diff         string
	ExperimentID int
	Left         File
	Right        File
}

// File contains the info of a File
type File struct {
	BlobID       string
	RepositoryID string
	CommitHash   string
	Path         string
	Content      string
	Hash         string
}

// Feature represents one name-value feature of file
type Feature struct {
	Name   string
	Weight float64
}

// Role represents the position of a app User
type Role string

// Value returns the string value of the Role
func (r Role) Value() (driver.Value, error) {
	if isValidRole(r) {
		return string(r), nil
	}

	return "", errors.New("invalid Role")
}

// Scan sets the Role with the passed string
func (r *Role) Scan(value interface{}) error {
	if v, ok := value.([]byte); ok && isValidRole(Role(string(v))) {
		*r = Role(v)
		return nil
	}

	return errors.New("can't scan a valid Role")
}

func isValidRole(r Role) bool {
	for _, role := range []Role{Worker, Requester} {
		if r == role {
			return true
		}
	}

	return false
}

const (
	// Requester is the role of a user that can review assignments, users, stats of experiments...
	Requester Role = "requester"
	// Worker is the role of a user that can answer assignments
	Worker Role = "worker"
)

// Answers lists the accepted answers
var Answers = map[string]string{
	"yes":   "yes",
	"maybe": "maybe",
	"no":    "no",
	"skip":  "skip",
}
