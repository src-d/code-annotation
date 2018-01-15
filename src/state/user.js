import api from '../api';

const initialState = {
  loggedIn: false,
  userId: null,
  username: null,
  avatarUrl: null,
};

export const LOG_IN = 'ca/user/LOG_IN';

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
    })
    .catch(e => {
      console.error(e);
    });

export const logOut = () => {};

export default reducer;
