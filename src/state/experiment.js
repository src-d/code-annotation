import { LOCATION_CHANGED, replace } from 'redux-little-router';
import api from '../api';
import { namedRoutes } from './routes';
import { add as addErrors } from './errors';
import { load as loadAssigments, undoneAssigment } from './assignments';

export const experimentId = 1; // hard-coded id for only experiment

/* reducer */

export const initialState = {
  loading: true,
  error: null,

  id: null,
  name: null,
};

export const LOAD = 'ca/experiment/LOAD';
export const LOAD_SUCCESS = 'ca/experiment/LOAD_SUCCESS';
export const LOAD_ERROR = 'ca/experiment/LOAD_ERROR';
export const SET_EXPERIMENT = 'ca/experiment/SET_EXPERIMENT';

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

    default:
      return state;
  }
};

export default reducer;

/* Actions */

export const load = expId => dispatch => {
  dispatch({ type: LOAD });
  return api
    .getExperiment(expId)
    .then(res =>
      dispatch({
        type: SET_EXPERIMENT,
        id: res.id,
        name: res.name,
      })
    )
    .then(() => dispatch({ type: LOAD_SUCCESS }))
    .catch(e => {
      dispatch(addErrors(e));
      dispatch({ type: LOAD_ERROR, error: e });
    });
};

/* Routing */

export const middleware = () => next => action => {
  if (action.type !== LOCATION_CHANGED) {
    return next(action);
  }

  const result = next(action);
  const { payload } = action;
  switch (payload.route) {
    case namedRoutes.experiment:
      return next(load(experimentId))
        .then(() => next(loadAssigments(experimentId)))
        .then(() => next(undoneAssigment(experimentId, replace)));
    default:
      return result;
  }
};
