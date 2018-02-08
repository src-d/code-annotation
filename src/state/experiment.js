import { createSelector } from 'reselect';
import { LOCATION_CHANGED, replace, push } from 'redux-little-router';
import api from '../api';
import { namedRoutes, makeUrl } from './routes';
import { add as addErrors } from './errors';

export const ANSWER_SIMILAR = 'yes';
export const ANSWER_MAYBE = 'maybe';
export const ANSWER_DIFFERENT = 'no';
export const ANSWER_SKIP = 'skip';

export const experimentId = 1; // hard-coded id for only experiment

/* reducer */

export const initialState = {
  loading: true,
  fileLoading: false,
  error: null,

  id: null,
  name: null,
  assignments: [],
  filePairs: {},

  currentAssigment: null,
  currentAssigmentStartTime: null,
};

export const LOAD = 'ca/experiment/LOAD';
export const LOAD_SUCCESS = 'ca/experiment/LOAD_SUCCESS';
export const LOAD_ERROR = 'ca/experiment/LOAD_ERROR';
export const SET_EXPERIMENT = 'ca/experiment/SET_EXPERIMENT';
export const SET_ASSIGNMENTS = 'ca/experiment/SET_ASSIGNMENTS';
export const LOAD_FILE_PAIR = 'ca/experiment/LOAD_FILE_PAIR';
export const SET_FILE_PAIR = 'ca/experiment/SET_FILE_PAIR';
export const SET_CURRENT_ASSIGMENT = 'ca/experiment/SET_CURRENT_ASSIGMENT';
export const MARK_ASSIGNMENT = 'ca/experiment/MARK_ASSIGNMENT';

const reducer = (state = initialState, action) => {
  switch (action.type) {
    case LOAD:
      return initialState;
    case LOAD_SUCCESS:
      return {
        ...state,
        loading: false,
      };
    case LOAD_ERROR:
      return {
        ...state,
        loading: false,
        fileLoading: false,
        error: action.error,
      };
    case SET_EXPERIMENT:
      return {
        ...state,
        id: action.id,
        name: action.name,
      };
    case SET_ASSIGNMENTS:
      return {
        ...state,
        assignments: action.assignments,
      };
    case LOAD_FILE_PAIR:
      return {
        ...state,
        fileLoading: true,
      };
    case SET_FILE_PAIR:
      return {
        ...state,
        fileLoading: false,
        filePairs: {
          ...state.filePairs,
          [action.id]: { diff: action.diff },
        },
      };
    case SET_CURRENT_ASSIGMENT:
      return {
        ...state,
        currentAssigment: action.assigment,
        currentAssigmentStartTime: new Date(),
      };
    case MARK_ASSIGNMENT:
      return {
        ...state,
        assignments: state.assignments.map(a => {
          if (a.id === action.id) {
            const previousDuration = a.duration || 0;
            const currentDuration =
              new Date() - state.currentAssigmentStartTime;
            return {
              ...a,
              answer: action.answer,
              duration: previousDuration + currentDuration,
            };
          }
          return a;
        }),
      };

    default:
      return state;
  }
};

export default reducer;

/* Actions */

const loadFilePairIfNeeded = id => (dispatch, getState) => {
  const { experiment } = getState();
  if (experiment.filePairs[id]) {
    return Promise.resolve();
  }

  dispatch({ type: LOAD_FILE_PAIR });

  return api.getFilePair(experiment.id, id).then(res =>
    dispatch({
      type: SET_FILE_PAIR,
      id: res.id,
      diff: res.diff,
    })
  );
};

const putAnswer = (expId, assignmentId, answer) => (dispatch, getState) => {
  const { experiment } = getState();
  const assigment = experiment.assignments.find(a => a.id === assignmentId);
  const answerPayload = {
    answer,
    duration: assigment.duration,
  };
  return api.putAnswer(expId, assignmentId, answerPayload).catch(e => {
    dispatch(addErrors(e));
    dispatch({ type: LOAD_ERROR, error: e });
  });
};

export const selectAssigment = id => (dispatch, getState) => {
  const { experiment } = getState();
  const assigment = experiment.assignments.find(a => a.id === +id);
  dispatch({
    type: SET_CURRENT_ASSIGMENT,
    assigment,
  });
  return dispatch(loadFilePairIfNeeded(assigment.pairId)).catch(e => {
    dispatch(addErrors(e));
    dispatch({ type: LOAD_ERROR, error: e });
  });
};

export const nextAssigment = (op = push) => (dispatch, getState) => {
  const { experiment } = getState();
  const noAnswer = experiment.assignments.find(a => a.answer === null);
  if (!noAnswer) {
    return dispatch(op(makeUrl('finish', { experiment: experimentId })));
  }

  return dispatch(
    op(makeUrl('question', { experiment: experimentId, question: noAnswer.id }))
  );
};

export const load = (expId, assigmentId) => dispatch => {
  dispatch({ type: LOAD });
  return api
    .getExperiment(expId)
    .then(res => {
      dispatch({
        type: SET_EXPERIMENT,
        id: res.id,
        name: res.name,
      });
      return api.getAssignments(res.id);
    })
    .then(res =>
      dispatch({
        type: SET_ASSIGNMENTS,
        assignments: res,
      })
    )
    .then(() => dispatch({ type: LOAD_SUCCESS }))
    .then(() => {
      if (!assigmentId) {
        return dispatch(nextAssigment(replace));
      }
      return dispatch(selectAssigment(assigmentId));
    })
    .catch(e => {
      dispatch(addErrors(e));
      dispatch({ type: LOAD_ERROR, error: e });
    });
};

export const mark = (assignmentId, answer) => dispatch => {
  dispatch({
    type: MARK_ASSIGNMENT,
    id: assignmentId,
    answer,
  });
  return dispatch(putAnswer(experimentId, assignmentId, answer)).then(() =>
    dispatch(nextAssigment())
  );
};

export const markCurrent = answer => (dispatch, getState) => {
  const { experiment } = getState();
  return dispatch(mark(experiment.currentAssigment.id, answer));
};

/* Routing */

export const middleware = store => next => action => {
  if (action.type !== LOCATION_CHANGED) {
    return next(action);
  }

  const result = next(action);
  const { payload } = action;
  const { experiment } = store.getState();
  switch (payload.route) {
    case namedRoutes.experiment:
      return next(load(experimentId));
    case namedRoutes.question:
      if (experiment.id === +payload.params.experiment) {
        return next(selectAssigment(+payload.params.question));
      }
      return next(load(experimentId, +payload.params.question));
    default:
      return result;
  }
};

/* Selectors */

const getAssignments = state => state.experiment.assignments;

export const getAssignmentsCount = state => state.experiment.assignments.length;

export const getIdenticalCount = createSelector(
  getAssignments,
  assignments => assignments.filter(a => a.answer === ANSWER_SIMILAR).length
);

export const getSimilarCount = createSelector(
  getAssignments,
  assignments => assignments.filter(a => a.answer === ANSWER_MAYBE).length
);

export const getDifferentCount = createSelector(
  getAssignments,
  assignments => assignments.filter(a => a.answer === ANSWER_DIFFERENT).length
);

export const getSkipCount = createSelector(
  getAssignments,
  assignments => assignments.filter(a => a.answer === ANSWER_SKIP).length
);

export const getProgressPercent = createSelector(
  getAssignmentsCount,
  getIdenticalCount,
  getSimilarCount,
  getDifferentCount,
  (total, indentical, similar, different) =>
    Math.round(100 / total * (indentical + similar + different))
);

export const getOverallTime = createSelector(getAssignments, assignments =>
  assignments.reduce((acc, a) => acc + (a.duration || 0), 0)
);

export const getAverageTime = createSelector(
  getOverallTime,
  getAssignments,
  (overallTime, assignments) =>
    Math.round(overallTime / assignments.filter(a => a.duration).length)
);

export const getCurrentFilePair = state => {
  const { currentAssigment, filePairs } = state.experiment;
  if (!currentAssigment) {
    return null;
  }
  return filePairs[currentAssigment.pairId];
};
