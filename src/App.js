import React, { Component } from 'react';
import { Fragment } from 'redux-little-router';
import Index from './pages/Index';
import Experiment from './pages/Experiment';

class App extends Component {
  render() {
    return (
      <Fragment forRoute="/">
        <div style={{ height: '100%' }}>
          <Fragment forRoute="/">
            <Index />
          </Fragment>
          <Fragment forRoute="/exp/:id">
            <Experiment />
          </Fragment>
          <Fragment forRoute="/exp/:id/:question">
            <Experiment />
          </Fragment>
          <Fragment forNoMatch>
            <div>not found</div>
          </Fragment>
        </div>
      </Fragment>
    );
  }
}

export default App;
