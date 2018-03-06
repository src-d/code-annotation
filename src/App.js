import React, { Component } from 'react';
import { Fragment } from 'redux-little-router';
import { Helmet } from 'react-helmet';
import { namedRoutes } from './state/routes';
import Errors from './components/Errors';
import Index from './pages/Index';
import Auth from './pages/Auth';
import Experiments from './pages/Experiments';
import Experiment from './pages/Experiment';
import Final from './pages/Final';
import Review from './pages/Review';
import Export from './pages/Export';
import Forbidden from './pages/Forbidden';

class App extends Component {
  render() {
    return (
      <Fragment forRoute="/">
        <div style={{ height: '100%' }}>
          <Helmet titleTemplate="%s | source{d} Code Annotation Tool" />
          <Errors />
          <Fragment forRoute={namedRoutes.index}>
            <Index />
          </Fragment>
          <Fragment forRoute={namedRoutes.auth}>
            <Auth />
          </Fragment>
          <Fragment forRoute={namedRoutes.dashboard}>
            <Experiments />
          </Fragment>
          <Fragment forRoute={namedRoutes.finish}>
            <Final />
          </Fragment>
          <Fragment forRoute={namedRoutes.review}>
            <Review />
          </Fragment>
          <Fragment forRoute={namedRoutes.question}>
            <Experiment />
          </Fragment>
          <Fragment forRoute={namedRoutes.experiment}>
            <Experiment />
          </Fragment>
          <Fragment forRoute={namedRoutes.export}>
            <Export />
          </Fragment>
          <Fragment forRoute={namedRoutes.forbidden}>
            <Forbidden />
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
