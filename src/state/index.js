import { combineReducers } from 'redux';
import user from './user';
import experiment from './experiment';

export default combineReducers({
  user,
  experiment,
});
