import api from '../api';

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
  return dispatch(loadFilePairIfNeeded(assigment.pairId)).then(() =>
    dispatch({
      type: SET_CURRENT_ASSIGMENT,
      assigment,
    })
  );
};

export const next = () => (dispatch, getState) => {
  const { experiment } = getState();
  const noAnswer = experiment.assignments.find(a => a.answer === null);
  if (!noAnswer) {
    // FIXME finish
    return Promise.resolve();
  }

  return dispatch(selectAssigment(noAnswer.id));
};

export const load = id => dispatch => {
  dispatch({ type: LOAD });
  return api
    .getExperiment(id)
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
    .then(() => dispatch(next()))
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
  return dispatch(next());
};

export const markCurrent = answer => (dispatch, getState) => {
  const { experiment } = getState();
  return dispatch(mark(experiment.currentAssigment.id, answer));
};
