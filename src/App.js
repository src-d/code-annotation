import React, { Component } from 'react';
import { Fragment } from 'redux-little-router';
import { namedRoutes } from './state/routes';
import Errors from './components/Errors';
import Index from './pages/Index';
import Experiment from './pages/Experiment';
import Final from './pages/Final';

class App extends Component {
  render() {
    return (
      <Fragment forRoute="/">
        <div style={{ height: '100%' }}>
          <Errors />
          <Fragment forRoute={namedRoutes.index}>
            <Index />
          </Fragment>
          <Fragment forRoute={namedRoutes.finish}>
            <Final />
          </Fragment>
          <Fragment forRoute={namedRoutes.question}>
            <Experiment />
          </Fragment>
          <Fragment forRoute={namedRoutes.experiment}>
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
