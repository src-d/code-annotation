import React, { Component } from 'react';
import { Button } from 'react-bootstrap';
import './Index.less';

class Index extends Component {
  render() {
    return (
      <div className="index-page">
        <div className="index-page__wrapper">
          <h1 className="index-page__header">
            Welcome! We&#39;re glad you made it this far.
          </h1>
          <p className="index-page__description">
            source{'{d}'} code annotation brings together state-of-art insights
            from<br />machine learning and user experience, for source code
            annotation.
          </p>

          <div className="index-page__actions">
            <Button
              href="/login"
              bsStyle="primary"
              className="index-page__github-button"
            >
              Sign in with Github
              <img
                src="/github.png"
                alt="github"
                className="index-page__github-icon"
              />
            </Button>
          </div>

          <div className="index-page__preview">
            <img src="/assigment.png" alt="assigment" width="650" />
          </div>
        </div>
      </div>
    );
  }
}

export default Index;
