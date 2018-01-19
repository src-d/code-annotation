import { LOCATION_CHANGED, replace, push } from 'redux-little-router';
import api from '../api';
import { add as addError } from './errors';
import TokenService from '../services/token';

const initialState = {
  loggedIn: false,
  userId: null,
  username: null,
  avatarUrl: null,
};

export const LOG_IN = 'ca/user/LOG_IN';
export const LOG_OUT = 'ca/user/LOG_OUT';

const reducer = (state = initialState, action) => {
  switch (action.type) {
    case LOG_IN:
      return {
        ...state,
        loggedIn: true,
        userId: action.userId,
        username: action.username,
        avatarUrl: action.avatarUrl,
      };
    case LOG_OUT: {
      return initialState;
    }

    default:
      return state;
  }
};

export const logIn = () => dispatch =>
  api.me().then(resp => {
    dispatch({
      type: LOG_IN,
      userId: resp.userId,
      username: resp.username,
      avatarUrl: resp.avatarURL,
    });
    dispatch(push('/exp/1'));
  });

export const logOut = () => dispatch => {
  TokenService.remove();
  dispatch({ type: LOG_OUT });
  return dispatch(push('/'));
};

export const authMiddleware = store => next => action => {
  if (action.type !== LOCATION_CHANGED) {
    return next(action);
  }
  const { query } = action.payload;
  let promise = Promise.resolve();
  if (query && query.token) {
    TokenService.set(query.token);
    promise = next(logIn());
  }

  return promise
    .then(() => {
      // redirect to / if try to access not public route wihout being logged in
      const { user } = store.getState();
      const { result } = action.payload;
      if (!user.loggedIn && result && !result.public) {
        return next(replace('/'));
      }
      return next(action);
    })
    .catch(e => next(addError(e)));
};

export default reducer;
