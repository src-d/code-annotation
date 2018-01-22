//TODO: connect with the DB
package repository

import (
	"fmt"

	"github.com/src-d/code-annotation/server/model"
)

var assignments []model.Assignment

var ErrNoAssignmentsInitialized = fmt.Errorf("No assignments initialized")

func GetExperimentById(ID int) (model.Experiment, error) {
	return model.Experiment{
		ID,
		fmt.Sprintf("Experiment#%d", ID),
		fmt.Sprintf("Description about experiment#%d", ID),
	}, nil
}

func GetAssignmentsFor(userID int, experimentID int) ([]model.Assignment, error) {
	if len(assignments) == 0 {
		return []model.Assignment{}, ErrNoAssignmentsInitialized
	}

	return assignments, nil
}

func CreateAssignmentsFor(userID int, experimentID int) ([]model.Assignment, error) {
	assignments = []model.Assignment{
		model.Assignment{ID: 1, UserID: userID, PairID: 1, ExperimentID: experimentID},
		model.Assignment{ID: 2, UserID: userID, PairID: 2, ExperimentID: experimentID},
		model.Assignment{ID: 3, UserID: userID, PairID: 1, ExperimentID: experimentID},
		model.Assignment{ID: 4, UserID: userID, PairID: 2, ExperimentID: experimentID},
		model.Assignment{ID: 5, UserID: userID, PairID: 1, ExperimentID: experimentID},
		model.Assignment{ID: 6, UserID: userID, PairID: 2, ExperimentID: experimentID},
		model.Assignment{ID: 7, UserID: userID, PairID: 1, ExperimentID: experimentID},
		model.Assignment{ID: 8, UserID: userID, PairID: 2, ExperimentID: experimentID},
	}

	return assignments, nil
}

func UpdateAssignment(assignmentId int, assignment model.Assignment) error {
	if _, ok := model.Answers[assignment.Answer]; !ok {
		return fmt.Errorf("Wrong answer provided: '%s'", assignment.Answer)
	}

	if assignmentId > len(assignments) {
		return fmt.Errorf("Assignment #%d does not exist", assignmentId)
	}

	assignments[assignmentId-1].Answer = assignment.Answer
	assignments[assignmentId-1].Duration = assignment.Duration
	return nil
}

func GetFilePairFor(pairID int) (model.FilePairs, error) {
	name := fmt.Sprintf("filePair-%d", pairID)
	return model.FilePairs{
		pairID,
		1,
		fmt.Sprintf("Diff of %s", name),
		model.File{"left", "ABC", fmt.Sprintf("%s left", name)},
		model.File{"right", "CBA", fmt.Sprintf("%s right", name)},
	}, nil
}
