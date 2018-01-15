import React, { Component } from 'react';
import { Button } from 'react-bootstrap';
import { connect } from 'react-redux';
import { logIn } from '../state/user';
import './Index.less';

class Index extends Component {
  render() {
    const { logIn: onClick } = this.props;
    return (
      <div className="index-page">
        <h1>Welcome! We&#39;re glad you make it this far.</h1>
        <p>
          source{'{d}'} code annotation brings together state-of-art insights
          from<br />machine learning and user experience, for source code
          annotation.
        </p>

        <Button onClick={onClick} bsStyle="primary">
          SIGN IN
        </Button>
      </div>
    );
  }
}

export default connect(undefined, { logIn })(Index);
