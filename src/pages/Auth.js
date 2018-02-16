import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Grid, Row, Col } from 'react-bootstrap';
import Loader from '../components/Loader';
import TokenService from '../services/token';
import { auth } from '../state/user';

class Auth extends Component {
  constructor(props) {
    super(props);

    this.state = {
      error: '',
    };
  }

  componentDidMount() {
    this.props.auth();
  }

  render() {
    const { error } = this.state;

    if (error) {
      return (
        <Grid>
          <Row className="ex-page__oops">
            <Col xs={12}>
              Oops.<br />Something went wrong.
            </Col>
          </Row>
          <Row style={{ paddingTop: '20px', paddingBottom: '20px' }}>
            <Col xs={12}>
              {error}
            </Col>
          </Row>
        </Grid>
      );
    }

    return (
      <Row className="ex-page__loader">
        <Col xs={12}>
          <Loader />
        </Col>
      </Row>
    );
  }
}

export default connect(undefined, { auth })(Auth);
