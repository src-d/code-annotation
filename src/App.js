import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Button } from 'react-bootstrap';
import './App.css';
import Noop from './components/Noop';

class App extends Component {
  render() {
    return (
      <div className="App">
        <header className="App-header">
          <h1 className="App-title">Welcome to React</h1>
        </header>
        <p className="App-intro">
          To get started, edit <code>src/App.js</code> and save to reload.
        </p>
        <Noop />
        <p>Test redux: {this.props.reduxValue}</p>
        <Button bsStyle="primary">bootstrap button</Button>
      </div>
    );
  }
}

const mapStateToProps = state => ({
  reduxValue: state.noop.value,
});

export default connect(mapStateToProps)(App);
