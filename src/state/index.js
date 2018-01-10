import { combineReducers } from 'redux';

const initialState = {
  value: 'hello from redux'
};

const noop = (state = initialState, action) => {
  return state;
};

export default combineReducers({
  noop
});
