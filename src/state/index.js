import { combineReducers } from 'redux';

const initialState = {
  value: 'hello from redux',
};

const noop = (state = initialState) => state;

export default combineReducers({
  noop,
});
