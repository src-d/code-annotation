//TODO: use "user.go" as a example, not this one ;)
package repository

import (
	"fmt"

	"github.com/src-d/code-annotation/server/model"
)

var assignments []*model.Assignment

var ErrNoAssignmentsInitialized = fmt.Errorf("No assignments initialized")

func GetExperimentByID(id int) (*model.Experiment, error) {
	return &model.Experiment{
		ID:          id,
		Name:        fmt.Sprintf("Experiment#%d", id),
		Description: fmt.Sprintf("Description about experiment#%d", id),
	}, nil
}

func GetAssignmentsFor(userID int, experimentID int) ([]*model.Assignment, error) {
	if len(assignments) == 0 {
		return []*model.Assignment{}, ErrNoAssignmentsInitialized
	}

	return assignments, nil
}

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
