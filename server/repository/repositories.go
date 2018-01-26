package repository

//TODO: use "user.go" as a example, not this one ;)

import (
	"fmt"

	"github.com/src-d/code-annotation/server/model"
)

var assignments []*model.Assignment

// ErrNoAssignmentsInitialized is the error returned when the Assignments of a User are requested
// for a given Experiment, but they have not been yet created
var ErrNoAssignmentsInitialized = fmt.Errorf("No assignments initialized")

// GetExperimentByID returns the Experiment identified by the passed ID
func GetExperimentByID(id int) (*model.Experiment, error) {
	return &model.Experiment{
		ID:          id,
		Name:        fmt.Sprintf("Experiment#%d", id),
		Description: fmt.Sprintf("Description about experiment#%d", id),
	}, nil
}

// GetAssignmentsFor returns the Assignments of a given user for a certain experiment,
// and returns an ErrNoAssignmentsInitialized if they does not yet exist
func GetAssignmentsFor(userID int, experimentID int) ([]*model.Assignment, error) {
	if len(assignments) == 0 {
		return []*model.Assignment{}, ErrNoAssignmentsInitialized
	}

	return assignments, nil
}

// CreateAssignmentsFor creates the assignments for the Experiment and User identified by the passed IDs
func CreateAssignmentsFor(userID int, experimentID int) ([]*model.Assignment, error) {
	assignments = []*model.Assignment{
		&model.Assignment{ID: 1, UserID: userID, PairID: 1, ExperimentID: experimentID},
		&model.Assignment{ID: 2, UserID: userID, PairID: 2, ExperimentID: experimentID},
		&model.Assignment{ID: 3, UserID: userID, PairID: 1, ExperimentID: experimentID},
		&model.Assignment{ID: 4, UserID: userID, PairID: 2, ExperimentID: experimentID},
		&model.Assignment{ID: 5, UserID: userID, PairID: 1, ExperimentID: experimentID},
		&model.Assignment{ID: 6, UserID: userID, PairID: 2, ExperimentID: experimentID},
		&model.Assignment{ID: 7, UserID: userID, PairID: 1, ExperimentID: experimentID},
		&model.Assignment{ID: 8, UserID: userID, PairID: 2, ExperimentID: experimentID},
	}

	return assignments, nil
}

// UpdateAssignment updates the Assignment identified by the passed ID, with the passed answer and duration
func UpdateAssignment(assignmentID int, answer string, duration int) error {
	if _, ok := model.Answers[answer]; !ok {
		return fmt.Errorf("Wrong answer provided: '%s'", answer)
	}

	if assignmentID > len(assignments) {
		return fmt.Errorf("Assignment #%d does not exist", assignmentID)
	}

	assignments[assignmentID-1].Answer = answer
	assignments[assignmentID-1].Duration = duration
	return nil
}

// GetFilePairFor returns the FilePairs identified by the passed ID
func GetFilePairFor(pairID int) (*model.FilePairs, error) {
	name := fmt.Sprintf("filePair-%d", pairID)
	return &model.FilePairs{
		ID:           pairID,
		ExperimentID: 1,
		Diff:         fmt.Sprintf("Diff of %s", name),
		Left:         model.File{Name: "left", Hash: "ABC", Content: fmt.Sprintf("%s left", name)},
		Right:        model.File{Name: "right", Hash: "CBA", Content: fmt.Sprintf("%s right", name)},
	}, nil
}
