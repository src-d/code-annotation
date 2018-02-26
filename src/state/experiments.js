import api from '../api';
import { add as addErrors } from './errors';

/* reducer */

export const initialState = {
  loading: true,
  error: null,

  list: [],
};

export const LOAD = 'ca/experiments/LOAD';
export const LOAD_SUCCESS = 'ca/experiments/LOAD_SUCCESS';
export const LOAD_ERROR = 'ca/experiments/LOAD_ERROR';
export const SET = 'ca/experiments/SET';

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
    case SET:
      return {
        ...state,
        list: action.data,
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
