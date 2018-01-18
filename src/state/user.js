import { LOCATION_CHANGED, replace, push } from 'redux-little-router';
import api from '../api';

const initialState = {
  loggedIn: false,
  userId: null,
  username: null,
  avatarUrl: null,
};

export const LOG_IN = 'ca/user/LOG_IN';

export const authMiddleware = store => next => action => {
  if (action.type !== LOCATION_CHANGED) {
    return next(action);
  }
  // redirect to / if try to access not public route wihout being logged in
  const { user } = store.getState();
  if (!user.loggedIn && !action.payload.result.public) {
    return next(replace('/'));
  }
  return next(action);
};

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

    default:
      return state;
  }
};

export const logIn = () => dispatch =>
  api
    .login()
    .then(resp => {
      dispatch({
        type: LOG_IN,
        userId: resp.userId,
        username: resp.username,
        avatarUrl: resp.avatarUrl,
      });
      dispatch(push('/exp/1'));
    })
    .catch(e => {
      console.error(e);
    });

export const logOut = () => {};

export default reducer;
