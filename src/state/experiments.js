import api from '../api';
import { add as addErrors } from './errors';

/* reducer */

export const UPLOAD_RES_NONE = 'ca/experiments/UPLOAD_RES_NONE';
export const UPLOAD_RES_SUCCESS = 'ca/experiments/UPLOAD_RES_SUCCESS';
export const UPLOAD_RES_FAILURE = 'ca/experiments/UPLOAD_RES_FAILURE';

export const initialState = {
  loading: true,
  error: null,
  createInProgress: false,
  updateInProgress: false,
  uploadInProgress: false,
  uploadResult: UPLOAD_RES_NONE,

  list: [],
};

export const LOAD = 'ca/experiments/LOAD';
export const LOAD_SUCCESS = 'ca/experiments/LOAD_SUCCESS';
export const LOAD_ERROR = 'ca/experiments/LOAD_ERROR';
export const SET = 'ca/experiments/SET';
export const CREATE = 'ca/experiments/CREATE';
export const CREATE_SUCCESS = 'ca/experiments/CREATE_SUCCESS';
export const CREATE_ERROR = 'ca/experiments/CREATE_ERROR';
export const UPDATE = 'ca/experiments/UPDATE';
export const UPDATE_SUCCESS = 'ca/experiments/UPDATE_SUCCESS';
export const UPDATE_ERROR = 'ca/experiments/UPDATE_ERROR';
export const UPLOAD = 'ca/experiments/UPLOAD';
export const UPLOAD_SUCCESS = 'ca/experiments/UPLOAD_SUCCESS';
export const UPLOAD_ERROR = 'ca/experiments/UPLOAD_ERROR';
export const UPLOAD_RESULT_RESET = 'ca/experiments/UPLOAD_RESULT_RESET';

const reducer = (state = initialState, action) => {
  switch (action.type) {
    case LOAD:
      return {
        ...state,
        loading: true,
        error: null,
        list: [],
      };
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
    case SET:
      return {
        ...state,
        list: action.data,
      };
    case CREATE:
      return {
        ...state,
        error: null,
        createInProgress: true,
      };
    case CREATE_SUCCESS:
      return {
        ...state,
        createInProgress: false,
      };
    case CREATE_ERROR:
      return {
        ...state,
        createInProgress: false,
        error: action.error,
      };
    case UPDATE:
      return {
        ...state,
        error: null,
        updateInProgress: true,
      };
    case UPDATE_SUCCESS:
      return {
        ...state,
        updateInProgress: false,
      };
    case UPDATE_ERROR:
      return {
        ...state,
        updateInProgress: false,
        error: action.error,
      };
    case UPLOAD:
      return {
        ...state,
        error: null,
        uploadInProgress: true,
        uploadResult: UPLOAD_RES_NONE,
      };
    case UPLOAD_SUCCESS:
      return {
        ...state,
        uploadInProgress: false,
        uploadResult: UPLOAD_RES_SUCCESS,
      };
    case UPLOAD_ERROR:
      return {
        ...state,
        uploadInProgress: false,
        uploadResult: UPLOAD_RES_FAILURE,
        error: action.error,
      };
    case UPLOAD_RESULT_RESET:
      return {
        ...state,
        uploadResult: UPLOAD_RES_NONE,
      };

    default:
      return state;
  }
};

export default reducer;

/* Actions */

export const load = () => dispatch => {
  dispatch({ type: LOAD });
  return api
    .getExperiments()
    .then(res =>
      dispatch({
        type: SET,
        data: res,
      })
    )
    .then(() => dispatch({ type: LOAD_SUCCESS }))
    .catch(e => {
      dispatch(addErrors(e));
      dispatch({ type: LOAD_ERROR, error: e });
    });
};

export const create = (name, description) => dispatch => {
  dispatch({ type: CREATE });
  return api
    .createExperiment(name, description)
    .then(res => {
      dispatch({ type: CREATE_SUCCESS });
      dispatch(load());
      return res.id;
    })
    .catch(e => {
      dispatch(addErrors(e));
      dispatch({ type: CREATE_ERROR, error: e });
      throw e;
    });
};

export const update = (id, name, description) => dispatch => {
  dispatch({ type: UPDATE });
  return api
    .updateExperiment(id, name, description)
    .then(() => {
      dispatch({ type: UPDATE_SUCCESS });
      dispatch(load());
    })
    .catch(e => {
      dispatch(addErrors(e));
      dispatch({ type: UPDATE_ERROR, error: e });
    });
};

export const uploadResultReset = () => dispatch => {
  dispatch({ type: UPLOAD_RESULT_RESET });
};

export const uploadFilePairs = (experimentId, file) => dispatch => {
  dispatch({ type: UPLOAD });
  return api
    .uploadFilePairs(experimentId, file)
    .then(() => {
      dispatch({ type: UPLOAD_SUCCESS });
    })
    .catch(e => {
      dispatch(addErrors(e));
      dispatch({ type: UPLOAD_ERROR, error: e });
    });
};
