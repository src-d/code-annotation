import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Grid, Row, Col } from 'react-bootstrap';
import PageHeader from '../components/PageHeader';
import AnnotationResults from '../components/AnnotationResults';
import {
  getIdenticalCount,
  getSimilarCount,
  getDifferentCount,
  getProgressPercent,
  getSkipCount,
  getAverageTime,
  getOverallTime,
} from '../state/assignments';

class Final extends Component {
  render() {
    const { user, expName, ...resultsProps } = this.props;

    return (
      <div className="final-page">
        <PageHeader {...user} />
        <Grid>
          <Row>
            <Col xs={12}>
              <h1 className="text-center">
                You did it great, {user.username}!
              </h1>
            </Col>
          </Row>
          <Row>
            <Col xs={12}>
              <p className="text-center">
                Here&#39;s an overview of {expName} experiment
              </p>
            </Col>
          </Row>
          <Row>
            <Col xs={6} xsOffset={3}>
              <AnnotationResults {...resultsProps} />
            </Col>
          </Row>
        </Grid>
      </div>
    );
  }
}

const mapStateToProps = state => {
  const { experiment } = state;
  return {
    user: state.user,
    expName: experiment.name,
    procent: getProgressPercent(state),
    identicalCount: getIdenticalCount(state),
    similarCount: getSimilarCount(state),
    differentCount: getDifferentCount(state),
    skipCount: getSkipCount(state),
    averageTime: getAverageTime(state),
    overallTime: getOverallTime(state),
  };
};

export default connect(mapStateToProps)(Final);
