package model

type User struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatarURL"`
	Role      Role
}

type Experiment struct {
	ID          int
	Name        string
	Description string
}

type Assignment struct {
	ID           int
	UserID       int
	PairID       int
	ExperimentID int
	Answer       string
	Duration     int
}

type FilePairs struct {
	ID           int
	ExperimentID int
	Diff         string
	Left         File
	Right        File
}

type File struct {
	Name    string
	Hash    string
	Content string
}

type Role string

const (
	Requester Role = "requester"
	Worker    Role = "worker"
)

var Answers = map[string]string{
	"yes":   "yes",
	"maybe": "maybe",
	"no":    "no",
	"skip":  "skip",
	"":      "",
}
