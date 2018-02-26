import React, { Component } from 'react';
import { Grid, Row, Col } from 'react-bootstrap';
import { connect } from 'react-redux';
import { push } from 'redux-little-router';
import { Helmet } from 'react-helmet';
import PageHeader from '../components/PageHeader';
import Loader from '../components/Loader';
import Breadcrumbs from '../components/Breadcrumbs';
import Progress from '../components/Experiment/Progress';
import Diff from '../components/Experiment/Diff';
import Selector from '../components/Experiment/Selector';
import Actions from '../components/Experiment/Actions';
import AdditionActions from '../components/Experiment/AdditionActions';
import {
  markCurrent,
  ANSWER_SIMILAR,
  ANSWER_MAYBE,
  ANSWER_DIFFERENT,
  ANSWER_SKIP,
  getProgressPercent,
  getCurrentFilePair,
} from '../state/assignments';
import { makeUrl } from '../state/routes';
import './Experiment.less';

class Experiment extends Component {
  render() {
    return (
      <div className="ex-page">
        <Helmet>
          <title>{this.props.name}</title>
        </Helmet>
        <PageHeader />
        {this.renderMain()}
      </div>
    );
  }

  renderMain() {
    const { error, loading, name, description, percent } = this.props;

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

    const breadcrumbsOptions = [{ name, link: '#' }];
    if (description) {
      breadcrumbsOptions.push({ name: description, link: '#' });
    }

    return (
      <Grid fluid className="ex-page__main">
        <Row className="ex-page__header">
          <Col xs={9} className="ex-page__info">
            <Breadcrumbs items={breadcrumbsOptions} />
          </Col>
          <Col xs={3} className="ex-page__progress">
            <Progress percent={percent} />
          </Col>
        </Row>
        {this.renderContent()}
      </Grid>
    );
  }

  renderContent() {
    const {
      expId,
      fileLoading,
      diffString,
      leftLoc,
      rightLoc,
      assignmentsOptions,
      currentAssigmentId,
      selectAssigmentId,
      markSimilar,
      markMaybe,
      markDifferent,
      skip,
      finish,
    } = this.props;

    if (fileLoading || !diffString) {
      return (
        <div className="ex-page__loader">
          <Loader />
        </div>
      );
    }

    return (
      <React.Fragment>
        <Row className="ex-page__content">
          <Col xs={12} className="ex-page__diff-col">
            <Diff
              diffString={diffString}
              leftLoc={leftLoc}
              rightLoc={rightLoc}
              className="ex-page__diff"
            />
          </Col>
        </Row>
        <Row className="ex-page__footer">
          <Col xs={3}>
            <Selector
              title="Previous"
              options={assignmentsOptions}
              value={currentAssigmentId}
              onChange={e => selectAssigmentId(expId, e.target.value)}
            />
          </Col>
          <Col xs={6} className="ex-page__actions">
            <Actions
              markSimilar={markSimilar}
              markMaybe={markMaybe}
              markDifferent={markDifferent}
            />
          </Col>
          <Col xs={3} className="ex-page__additional-actions">
            <AdditionActions skip={skip} finish={() => finish(expId)} />
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

const mapStateToProps = state => {
  const { experiment, assignments } = state;
  const { error, loading, id, name, description } = experiment;
  const { list, currentAssigment } = assignments;

  const filePair = getCurrentFilePair(state);
  const diff = filePair ? filePair.diff : null;

  const assignmentsOptions = list.map((a, i) => {
    const status = a.answer ? ` (${a.answer})` : '';
    return { value: a.id, name: `(${i + 1})${status}` };
  });

  return {
    error,
    loading,
    expId: id,
    name,
    description,
    percent: getProgressPercent(state),
    diffString: diff,
    leftLoc: filePair ? filePair.leftLoc : 0,
    rightLoc: filePair ? filePair.rightLoc : 0,
    currentAssigmentId: currentAssigment ? currentAssigment.id : null,
    assignmentsOptions,
  };
};

const mapDispatchToProps = dispatch => ({
  markSimilar: () => dispatch(markCurrent(ANSWER_SIMILAR)),
  markMaybe: () => dispatch(markCurrent(ANSWER_MAYBE)),
  markDifferent: () => dispatch(markCurrent(ANSWER_DIFFERENT)),
  skip: () => dispatch(markCurrent(ANSWER_SKIP)),
  selectAssigmentId: (expId, id) =>
    dispatch(push(makeUrl('question', { experiment: expId, question: id }))),
  finish: expId => dispatch(push(makeUrl('finish', { experiment: expId }))),
});

export default connect(mapStateToProps, mapDispatchToProps)(Experiment);
