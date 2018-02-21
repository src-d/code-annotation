import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Grid, Row, Col } from 'react-bootstrap';
import { Helmet } from 'react-helmet';
import PageHeader from '../components/PageHeader';
import ExperimentsList from '../components/ExperimentsList';

class Final extends Component {
  render() {
    const { user, experiments } = this.props;

    return (
      <div className="experiments-page">
        <Helmet>
          <title>Experiments</title>
        </Helmet>
        <PageHeader />
        <Grid>
          <Row>
            <Col xs={12}>
              <h1 className="text-center">
                You look great today, {user.username}!
              </h1>
            </Col>
          </Row>
          <Row>
            <Col xs={12}>
              <p className="text-center">
                Here&#39;s your experiments dashboard, go ahead and make our
                day.
              </p>
            </Col>
          </Row>
          <Row>
            <Col xs={8} xsOffset={2} style={{ paddingTop: '40px' }}>
              <ExperimentsList experiments={experiments} />
            </Col>
          </Row>
        </Grid>
      </div>
    );
  }
}

const mapStateToProps = state => ({
  user: state.user,
  experiments: [
    { id: 1, name: 'Lee Morgan Experiment', procent: 0 },
    { id: 2, name: 'Chet Baker Experiment', procent: 36 },
    { id: 2, name: 'Cliff Brown Experiment', procent: 100 },
  ],
});

export default connect(mapStateToProps)(Final);
