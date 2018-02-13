import { createSelector } from 'reselect';
import { add as addErrors } from './errors';
import api from '../api';

export const initialState = {
  loading: false,
  error: null,
  features: [],
};

export const LOAD = 'ca/features/LOAD';
export const LOAD_SUCCESS = 'ca/features/LOAD_SUCCESS';

const reducer = (state = initialState, action) => {
  switch (action.type) {
    case LOAD:
      return {
        ...state,
        loading: true,
        error: null,
        features: [],
      };
    case LOAD_SUCCESS:
      return {
        ...state,
        loading: false,
        features: action.features,
      };
    default:
      return state;
  }
};

// Actions

export const load = (blobIdA, blobIdB) => dispatch => {
  dispatch({ type: LOAD });

  return Promise.all([api.getFeatures(blobIdA), api.getFeatures(blobIdB)])
    .then(([featuresA, featuresB]) => {
      // collect names from both responses
      const names = featuresA.map(f => f.name);
      featuresB.forEach(f => {
        if (!names.includes(f.name)) {
          names.push(f.name);
        }
      });

      // convert lists to maps
      const mapA = featuresA.reduce((acc, f) => {
        acc[f.name] = f.weight;
        return acc;
      }, {});
      const mapB = featuresB.reduce((acc, f) => {
        acc[f.name] = f.weight;
        return acc;
      }, {});

      // merge features by name
      const features = names.map(name => ({
        name,
        weightA: mapA[name] || 0,
        weightB: mapB[name] || 0,
      }));

      return dispatch({ type: LOAD_SUCCESS, features });
    })
    .catch(err => {
      dispatch(addErrors(err));
    });
};

// Selectors

export const getFeatures = state => state.features.features;

export const mostSimilar = createSelector(getFeatures, features => {
  // copy because sort mutates array
  const list = [...features];
  return list.sort((a, b) => {
    const deltaA = Math.abs(a.weightA - a.weightB);
    const deltaB = Math.abs(b.weightA - b.weightB);
    if (deltaA > deltaB) {
      return 1;
    }
    if (deltaB > deltaA) {
      return -1;
    }
    if (a.weightA > b.weightA) {
      return -1;
    }
    if (b.weightA > a.weightA) {
      return 1;
    }
    return 0;
  });
});

export const leastSimilar = createSelector(mostSimilar, features => {
  // reverse mutates
  const list = [...features];
  return list.reverse();
});

export default reducer;
