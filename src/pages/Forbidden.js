import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Grid, Row, Col } from 'react-bootstrap';
import PageHeader from '../components/PageHeader';

class Forbidden extends Component {
  render() {
    return (
      <div className="forbidden-page">
        <PageHeader {...this.props.user} />
        <Grid>
          <Row style={{ paddingTop: '20px', paddingBottom: '20px' }}>
            <Col xs={12}>
              <p>
                You tried to perform an action you are not currently authorized
                to do.
              </p>
            </Col>
          </Row>
        </Grid>
      </div>
    );
  }
}

const mapStateToProps = state => ({
  user: state.user,
});

export default connect(mapStateToProps)(Forbidden);
