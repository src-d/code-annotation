import { createStore } from 'redux';
import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import rootReducer from './state';

const store = createStore(rootReducer);

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(<App store={store} />, div);
});
