import { createSelector } from 'reselect';
import { LOCATION_CHANGED, push } from 'redux-little-router';
import api from '../api';
import { namedRoutes, makeUrl } from './routes';
import { add as addErrors } from './errors';
import { load as expLoad } from './experiment';
import { loadFilePairIfNeeded } from './filePairs';

export const ANSWER_SIMILAR = 'yes';
export const ANSWER_MAYBE = 'maybe';
export const ANSWER_DIFFERENT = 'no';
export const ANSWER_SKIP = 'skip';

export const experimentId = 1; // hard-coded id for only experiment

/* reducer */

export const initialState = {
  loading: true,
  error: null,

  list: [],

  currentAssigment: null,
  currentAssigmentStartTime: null,
};

export const LOAD = 'ca/assignments/LOAD';
export const LOAD_SUCCESS = 'ca/assignments/LOAD_SUCCESS';
export const LOAD_ERROR = 'ca/assignments/LOAD_ERROR';
export const SET_ASSIGNMENTS = 'ca/assignments/SET_ASSIGNMENTS';
export const SET_CURRENT_ASSIGMENT = 'ca/assignments/SET_CURRENT_ASSIGMENT';
export const MARK_ASSIGNMENT = 'ca/assignments/MARK_ASSIGNMENT';

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
        error: action.error,
      };
    case SET_ASSIGNMENTS:
      return {
        ...state,
        list: action.assignments,
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
        list: state.list.map(a => {
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

const putAnswer = (expId, assignmentId, answer) => (dispatch, getState) => {
  const { assignments } = getState();
  const assigment = assignments.list.find(a => a.id === assignmentId);
  const answerPayload = {
    answer,
    duration: assigment.duration,
  };
  return api.putAnswer(expId, assignmentId, answerPayload);
};

export const selectAssigment = (expId, id) => (dispatch, getState) => {
  const { assignments } = getState();
  const assigment = assignments.list.find(a => a.id === +id);
  dispatch({
    type: SET_CURRENT_ASSIGMENT,
    assigment,
  });
  return dispatch(loadFilePairIfNeeded(expId, assigment.pairId)).catch(e => {
    dispatch(addErrors(e));
    dispatch({ type: LOAD_ERROR, error: e });
  });
};

const goToQuestion = (expId, question, op = push) => dispatch => {
  if (!question) {
    return dispatch(op(makeUrl('finish', { experiment: expId })));
  }

  return dispatch(
    op(makeUrl('question', { experiment: expId, question: question.id }))
  );
};

export const undoneAssigment = (expId, op = push) => (dispatch, getState) => {
  const { assignments } = getState();
  const noAnswer = assignments.list.find(a => a.answer === null);
  return dispatch(goToQuestion(expId, noAnswer, op));
};

export const nextAssigment = (expId, op = push) => (dispatch, getState) => {
  const { assignments } = getState();
  let currentIdx;
  if (assignments.currentAssigment) {
    currentIdx = assignments.list.findIndex(
      a => a.id === assignments.currentAssigment.id
    );
  } else {
    currentIdx = -1;
  }

  if (currentIdx >= assignments.list.length - 1) {
    return dispatch(undoneAssigment(expId, op));
  }

  const question = assignments.list[currentIdx + 1];
  return dispatch(goToQuestion(expId, question, op));
};

export const load = expId => dispatch => {
  dispatch({ type: LOAD });
  return api
    .getAssignments(expId)
    .then(res =>
      dispatch({
        type: SET_ASSIGNMENTS,
        assignments: res,
      })
    )
    .then(() => dispatch({ type: LOAD_SUCCESS }))
    .catch(e => {
      dispatch(addErrors(e));
      dispatch({ type: LOAD_ERROR, error: e });
    });
};

export const mark = (expId, assignmentId, answer) => dispatch => {
  dispatch({
    type: MARK_ASSIGNMENT,
    id: assignmentId,
    answer,
  });
  return dispatch(putAnswer(expId, assignmentId, answer))
    .then(() => dispatch(nextAssigment(expId)))
    .catch(e => {
      dispatch(addErrors(e));
      dispatch({ type: LOAD_ERROR, error: e });
    });
};

export const markCurrent = answer => (dispatch, getState) => {
  const { experiment, assignments } = getState();
  return dispatch(mark(experiment.id, assignments.currentAssigment.id, answer));
};

/* Routing */

export const middleware = store => next => action => {
  if (action.type !== LOCATION_CHANGED) {
    return next(action);
  }

  const result = next(action);
  const { payload } = action;
  const { experiment } = store.getState();
  let promise;
  switch (payload.route) {
    case namedRoutes.question:
      promise = Promise.resolve();
      if (experiment.id !== +payload.params.experiment) {
        promise = next(expLoad(experimentId));
      }
      return promise
        .then(() => next(load(experimentId)))
        .then(() =>
          next(selectAssigment(experimentId, +payload.params.question))
        );
    case namedRoutes.finish:
      promise = Promise.resolve();
      if (experiment.id !== +payload.params.experiment) {
        promise = next(expLoad(experimentId));
      }
      return promise.then(() => next(load(experimentId)));
    default:
      return result;
  }
};

/* Selectors */

const getAssignments = state => state.assignments.list;

export const getAssignmentsCount = state => state.assignments.list.length;

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
  (total, indentical, similar, different) => {
    const completed = indentical + similar + different;
    if (!total || !completed) {
      return 0;
    }
    return Math.round(100 / total * completed);
  }
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
  const { currentAssigment } = state.assignments;
  const { pairs } = state.filePairs;
  if (!currentAssigment) {
    return null;
  }
  return pairs[currentAssigment.pairId];
};
