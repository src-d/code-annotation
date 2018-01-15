import React, { Component } from 'react';
import { connect } from 'react-redux';
import Index from './pages/Index';
import Experiment from './pages/Experiment';

class App extends Component {
  render() {
    const { loggedIn } = this.props;
    if (!loggedIn) {
      return <Index />;
    }
    return <Experiment />;
  }
}

const mapStateToProps = state => ({
  loggedIn: state.user.loggedIn,
});

export default connect(mapStateToProps)(App);
