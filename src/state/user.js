import { LOCATION_CHANGED, replace, push } from 'redux-little-router';
import api from '../api';
import { add as addError } from './errors';
import { makeUrl } from './routes';
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
        role: action.role,
      };
    case LOG_OUT: {
      return initialState;
    }

    default:
      return state;
  }
};

export const logOut = () => dispatch => {
  TokenService.remove();
  dispatch({ type: LOG_OUT });
  return dispatch(push('/'));
};

export const logIn = () => dispatch =>
  api
    .me()
    .then(resp => {
      dispatch({
        type: LOG_IN,
        userId: resp.id,
        username: resp.username,
        avatarUrl: resp.avatarURL,
        role: resp.role,
      });
    })
    .catch(e => {
      dispatch(addError(e));
      return dispatch(logOut());
    });

export const auth = () => (dispatch, getState) => {
  const { router } = getState();
  return api
    .auth(router.search)
    .then(data => {
      TokenService.set(data.token);
      return dispatch(logIn());
    })
    .then(() => dispatch(push('/dashboard')))
    .catch(e => {
      dispatch(addError(e));
      return dispatch(logOut());
    });
};

export const authMiddleware = store => next => action => {
  if (action.type !== LOCATION_CHANGED) {
    return next(action);
  }

  // redirect to / if try to access not public route wihout being logged in
  const { user } = store.getState();
  const { result } = action.payload;
  if (!user.loggedIn && result && !result.public) {
    return next(replace('/'));
  }

  // redirect user from index page to experiment if user it authorized already
  if (user.loggedIn && result && result.name === 'index') {
    return next(push(makeUrl('dashboard')));
  }

  // hide pages that are meant only for users with the requester role
  if (user.role !== 'requester' && result && result.restrictReviewer) {
    return next(replace('/forbidden'));
  }

  return next(action);
};

export default reducer;
