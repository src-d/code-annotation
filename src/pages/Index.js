import React, { Component } from 'react';
import { Grid, Row, Col, Button } from 'react-bootstrap';
import { Helmet } from 'react-helmet';
import { serverUrl } from '../api';
import './Index.less';

class Index extends Component {
  render() {
    return (
      <Grid fluid className="index-page">
        <Helmet>
          <title>Dashboard</title>
        </Helmet>
        <Row>
          <Col xs={12}>
            <h1 className="index-page__header">
              Welcome! We&#39;re glad you made it this far.
            </h1>
          </Col>
        </Row>
        <Row>
          <Col xs={12}>
            <p className="index-page__description">
              source{'{d}'} code annotation brings together state-of-art
              insights from<br />machine learning and user experience, for
              source code annotation.
            </p>
          </Col>
        </Row>
        <Row>
          <Col xs={12} className="index-page__actions">
            <Button
              href={`${serverUrl}/login`}
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
          </Col>
        </Row>
        <Row>
          <Col xs={12} className="index-page__preview">
            <img src="/assigment.png" alt="assigment" width="650" />
          </Col>
        </Row>
      </Grid>
    );
  }
}

export default Index;
