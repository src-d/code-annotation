package model

// User of the application; can be Requester or Workers
type User struct {
	ID        int
	Login     string
	Username  string
	AvatarURL string
	Role      Role
}

// Experiment groups a certain amount of FilePairs
type Experiment struct {
	ID          int
	Name        string
	Description string
}

// Assignment tracks the answer of a worker to a given FilePairs of an Experiment
type Assignment struct {
	ID           int
	UserID       int
	PairID       int
	ExperimentID int
	Answer       string
	Duration     int
}

// FilePairs are the answers that needs to be responsed
type FilePairs struct {
	ID           int
	ExperimentID int
	Diff         string
	Left         File
	Right        File
}

// File contains the info of a File
type File struct {
	Name    string
	Hash    string
	Content string
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
	"":      "",
}
