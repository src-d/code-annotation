import 'bootstrap/dist/css/bootstrap.css';

import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import { compose, createStore, applyMiddleware } from 'redux';
import { initializeCurrentLocation } from 'redux-little-router';
import { composeWithDevTools } from 'redux-devtools-extension';
import thunk from 'redux-thunk';
import rootReducer, { middlewares } from './state';
import { enhancer as routerEnhancer } from './state/routes';
import { logIn } from './state/user';
import TokenService from './services/token';
import './override.less';
import './index.less';
import App from './App';

const store = createStore(
  rootReducer,
  composeWithDevTools(
    compose(routerEnhancer, applyMiddleware(...[...middlewares, thunk]))
  )
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
