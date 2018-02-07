package model

import "database/sql"

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

// Role represents the position of a app User
type Role string

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
