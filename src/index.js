import 'bootstrap/dist/css/bootstrap.css';

import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import { compose, createStore, applyMiddleware } from 'redux';
import { initializeCurrentLocation } from 'redux-little-router';
import { composeWithDevTools } from 'redux-devtools-extension';
import thunk from 'redux-thunk';
import rootReducer, { middlewares as stateMiddlewares } from './state';
import { enhancer as routerEnhancer } from './state/routes';
import { logIn } from './state/user';
import TokenService from './services/token';
import { load as loadGA, middleware as gaMiddleware } from './services/ga';
import './override.less';
import './index.less';
import App from './App';

const middlewares = [...stateMiddlewares, thunk];

if (window.cat.GA_TRACKING_ID) {
  loadGA(window.cat.GA_TRACKING_ID);
  middlewares.push(gaMiddleware);
}

const store = createStore(
  rootReducer,
  window.cat.initialState || {},
  composeWithDevTools(compose(routerEnhancer, applyMiddleware(...middlewares)))
);

let promise = Promise.resolve();
if (TokenService.exists()) {
  promise = store.dispatch(logIn());
}

promise.catch().then(() => {
  const initialLocation = store.getState().router;
  if (initialLocation) {
    store.dispatch(initializeCurrentLocation(initialLocation));
  }

  ReactDOM.render(
    <Provider store={store}>
      <App />
    </Provider>,
    document.getElementById('root')
  );
});
