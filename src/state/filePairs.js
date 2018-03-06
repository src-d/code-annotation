import { LOCATION_CHANGED, replace } from 'redux-little-router';
import api from '../api';
import { add as addErrors } from './errors';
import { namedRoutes, makeUrl } from './routes';
import { load as expLoad } from './experiment';
import { load as featuresLoad } from './features';

/* reducer */

export const initialState = {
  loading: true,
  list: [],
  pairs: {},

  currentPairId: null,
  currentAnnotations: {},
};

export const LOAD = 'ca/filePairs/LOAD';
export const LOAD_SUCCESS = 'ca/filePairs/LOAD_SUCCESS';
export const LOAD_FILE_PAIR = 'ca/filePairs/LOAD_FILE_PAIR';
export const LOAD_ANNOTATIONS = 'ca/filePairs/LOAD_ANNOTATIONS';
export const SET_ANNOTATIONS = 'ca/filePairs/SET_ANNOTATIONS';
export const SET_FILE_PAIR = 'ca/filePairs/SET_FILE_PAIR';
export const SET_CURRENT_PAIR = 'ca/filePairs/SET_CURRENT_PAIR';

const reducer = (state = initialState, action) => {
  switch (action.type) {
    case LOAD:
      return initialState;
    case LOAD_SUCCESS:
      return {
        ...state,
        loading: false,
        list: action.pairs,
      };
    case LOAD_FILE_PAIR:
      return {
        ...state,
        loading: true,
      };
    case SET_FILE_PAIR:
      return {
        ...state,
        loading: false,
        pairs: {
          ...state.pairs,
          [action.id]: {
            diff: action.diff,
            score: action.score,
            leftBlobId: action.leftBlobId,
            rightBlobId: action.rightBlobId,
            leftLoc: action.leftLoc,
            rightLoc: action.rightLoc,
          },
        },
      };
    case LOAD_ANNOTATIONS:
      return {
        ...state,
        currentAnnotations: {},
      };
    case SET_ANNOTATIONS:
      return {
        ...state,
        currentAnnotations: action.data,
      };
    case SET_CURRENT_PAIR:
      return {
        ...state,
        currentPairId: action.id,
      };

    default:
      return state;
  }
};

export default reducer;

/* Actions */

export const load = expId => dispatch => {
  dispatch({ type: LOAD });

  return api.getFilePairs(expId).then(pairs =>
    dispatch({
      type: LOAD_SUCCESS,
      pairs,
    })
  );
};

export const loadFilePairIfNeeded = (expId, id) => (dispatch, getState) => {
  const { filePairs } = getState();
  if (filePairs.pairs[id]) {
    return Promise.resolve();
  }

  dispatch({ type: LOAD_FILE_PAIR });

  return api.getFilePair(expId, id).then(res =>
    dispatch({
      type: SET_FILE_PAIR,
      ...res,
    })
  );
};

export const selectPair = (expId, id) => dispatch => {
  dispatch({
    type: SET_CURRENT_PAIR,
    id,
  });

  return dispatch(loadFilePairIfNeeded(expId, id)).catch(e => {
    dispatch(addErrors(e));
  });
};

export const loadAnnotations = (expId, id) => dispatch => {
  dispatch({ type: LOAD_ANNOTATIONS });
  return api
    .getFilePairAnnotations(expId, id)
    .then(res =>
      dispatch({
        type: SET_ANNOTATIONS,
        data: res,
      })
    )
    .catch(e => {
      dispatch(addErrors(e));
    });
};

/* Selectors */

export const getCurrentFilePair = state => {
  const { currentPairId, pairs } = state.filePairs;
  if (!currentPairId) {
    return null;
  }
  return pairs[currentPairId];
};

/* Routing */

export const middleware = store => next => action => {
  if (action.type !== LOCATION_CHANGED) {
    return next(action);
  }

  const result = next(action);
  const { payload } = action;
  const { experiment } = store.getState();
  let expIdParam;
  let promise;
  switch (payload.route) {
    case namedRoutes.review:
      expIdParam = +payload.params.experiment;
      return next(expLoad(expIdParam)).then(() =>
        next(load(expIdParam)).then(() => {
          const { filePairs } = store.getState();
          // FIXME(max) can be empty, but I really really want to fix it in separate PR
          // I'm sure it's not the only place where we handle empty responses incorrectly
          const pairId = filePairs.list[0].id;
          return next(
            replace(
              makeUrl('reviewPair', {
                experiment: expIdParam,
                pair: pairId,
              })
            )
          );
        })
      );
    case namedRoutes.reviewPair:
      expIdParam = +payload.params.experiment;
      promise = Promise.resolve();
      if (experiment.id !== expIdParam) {
        promise = next(expLoad(expIdParam)).then(() => next(load(expIdParam)));
      }
      return promise
        .then(() => next(selectPair(expIdParam, +payload.params.pair)))
        .then(() => next(loadAnnotations(expIdParam, +payload.params.pair)))
        .then(() => {
          const pair = getCurrentFilePair(store.getState());
          return pair && next(featuresLoad(pair.leftBlobId, pair.rightBlobId));
        });
    default:
      return result;
  }
};
