import { combineReducers } from 'redux';
import user, { authMiddleware } from './user';
import experiments from './experiments';
import experiment, { middleware as expMiddleware } from './experiment';
import assignments, {
  middleware as assignmentsMiddleware,
} from './assignments';
import filePairs, { middleware as filePairsMiddleware } from './filePairs';
import features from './features';
import router, { middleware as routerMiddleware } from './routes';
import errors from './errors';

export const middlewares = [
  authMiddleware,
  expMiddleware,
  assignmentsMiddleware,
  filePairsMiddleware,
  routerMiddleware,
];

export default combineReducers({
  router,
  errors,
  user,
  experiments,
  experiment,
  assignments,
  filePairs,
  features,
});
