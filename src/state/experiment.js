import { LOCATION_CHANGED, replace, push } from 'redux-little-router';
import api from '../api';
import { namedRoutes, makeUrl } from './routes';

export const ANSWER_SIMILAR = 'Yes';
export const ANSWER_MAYBE = 'Maybe';
export const ANSWER_DIFFERENT = 'No';
export const ANSWER_SKIP = 'Skip';

export const experimentId = 1; // hard-coded id for only experiment

/* reducer */

const initialState = {
  loading: true,
  fileLoading: false,
  error: null,

  id: null,
  name: null,
  assignments: [],
  filePairs: {},

  currentAssigment: null,
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
      };
    case MARK_ASSIGNMENT:
      return {
        ...state,
        assignments: state.assignments.map(a => {
          if (a.id === action.id) {
            return { ...a, answer: action.answer };
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

export const loadFilePairIfNeeded = id => (dispatch, getState) => {
  const { experiment } = getState();
  if (experiment.filePairs[id]) {
    return Promise.resolve();
  }

  dispatch({ type: LOAD_FILE_PAIR });

  return api.getFilePair(id).then(res =>
    dispatch({
      type: SET_FILE_PAIR,
      id: res.id,
      diff: res.diff,
    })
  );
};

export const selectAssigment = id => (dispatch, getState) => {
  const { experiment } = getState();
  const assigment = experiment.assignments.find(a => a.id === +id);
  dispatch({
    type: SET_CURRENT_ASSIGMENT,
    assigment,
  });
  return dispatch(loadFilePairIfNeeded(assigment.pairId));
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
    .then(res => {
      dispatch({
        type: SET_ASSIGNMENTS,
        assignments: res,
      });
    })
    .then(() => {
      if (!assigmentId) {
        return dispatch(nextAssigment(replace));
      }
      return dispatch(selectAssigment(assigmentId));
    })
    .then(() => {
      dispatch({ type: LOAD_SUCCESS });
    })
    .catch(e => {
      console.error(e);
      dispatch({ type: LOAD_ERROR, error: e });
    });
};

export const mark = (assignmentId, answer) => dispatch => {
  dispatch({
    type: MARK_ASSIGNMENT,
    id: assignmentId,
    answer,
  });
  return dispatch(nextAssigment());
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

/* Selectors (we might need to use reselect in future) */

export const getAssignmentsCount = state => state.experiment.assignments.length;

export const getSimilarCount = state =>
  state.experiment.assignments.filter(a => a.answer === ANSWER_SIMILAR).length;

export const getDifferentCount = state =>
  state.experiment.assignments.filter(a => a.answer === ANSWER_DIFFERENT)
    .length;

export const getSkipCount = state =>
  state.experiment.assignments.filter(a => a.answer === ANSWER_SKIP).length;

export const getProgressPercent = state =>
  100 /
  getAssignmentsCount(state) *
  (getSimilarCount(state) + getDifferentCount(state));
