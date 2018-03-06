import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Grid, Row, Col } from 'react-bootstrap';
import Page from './Page';
import Loader from '../components/Loader';
import ExperimentsList from '../components/ExperimentsList';
import { load as experimentsLoad } from '../state/experiments';

class Experiments extends Component {
  componentDidMount() {
    this.props.load();
  }

  render() {
    return <Grid>{this.renderContent()}</Grid>;
  }

  renderContent() {
    const { user, experiments } = this.props;
    const { error, loading, list } = experiments;

    if (error) {
      return (
        <Row className="ex-page__oops">
          <Col xs={12}>
            Oops.<br />Something went wrong.
          </Col>
        </Row>
      );
    }

    if (loading) {
      return (
        <Row className="ex-page__loader">
          <Col xs={12}>
            <Loader />
          </Col>
        </Row>
      );
    }

    return (
      <React.Fragment>
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
              Here&#39;s your experiments dashboard, go ahead and make our day.
            </p>
          </Col>
        </Row>
        <Row>
          <Col xs={8} xsOffset={2} style={{ paddingTop: '40px' }}>
            <ExperimentsList experiments={list} />
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

const mapStateToProps = state => ({
  user: state.user,
  experiments: state.experiments,
});

export default connect(mapStateToProps, { load: experimentsLoad })(
  Page(Experiments, {
    className: 'experiments-page',
    titleFn: () => 'Experiments',
    showHeader: true,
  })
);
