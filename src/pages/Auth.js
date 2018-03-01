import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Grid, Row, Col } from 'react-bootstrap';
import { Helmet } from 'react-helmet';
import PageHeader from '../components/PageHeader';
import Loader from '../components/Loader';
import { auth } from '../state/user';
import './Auth.less';

class Auth extends Component {
  componentDidMount() {
    this.props.auth();
  }

  render() {
    return (
      <div className="auth-page">
        <Helmet>
          <title>Authorization</title>
        </Helmet>
        <PageHeader />
        <Grid>
          <Row>
            <Col xs={12} className="auth-page__loader">
              <Loader />
            </Col>
          </Row>
        </Grid>
      </div>
    );
  }
}

export default connect(undefined, { auth })(Auth);
