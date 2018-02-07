import React, { Component } from 'react';
import { Grid, Row, Col } from 'react-bootstrap';
import { connect } from 'react-redux';
import { push } from 'redux-little-router';
import PageHeader from '../components/PageHeader';
import Loader from '../components/Loader';
import Progress from '../components/Experiment/Progress';
import Diff from '../components/Experiment/Diff';
import Selector from '../components/Experiment/Selector';
import Actions from '../components/Experiment/Actions';
import AdditionActions from '../components/Experiment/AdditionActions';
import {
  load,
  markCurrent,
  ANSWER_SIMILAR,
  ANSWER_MAYBE,
  ANSWER_DIFFERENT,
  ANSWER_SKIP,
  experimentId,
  getProgressPercent,
  getCurrentFilePair,
} from '../state/experiment';
import { makeUrl } from '../state/routes';
import './Experiment.less';

class Experiment extends Component {
  render() {
    const { user } = this.props;

    return (
      <div className="ex-page">
        <PageHeader {...user} />
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

    return (
      <Grid fluid className="ex-page__main">
        <Row className="ex-page__header">
          <Col xs={9} className="ex-page__info">
            <span className="ex-page__name">{name}</span>
            <span className="ex-page__description">{description}asd</span>
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
      fileLoading,
      diffString,
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
            <Diff diffString={diffString} className="ex-page__diff" />
          </Col>
        </Row>
        <Row className="ex-page__footer">
          <Col xs={3}>
            <Selector
              options={assignmentsOptions}
              value={currentAssigmentId}
              select={selectAssigmentId}
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
            <AdditionActions skip={skip} finish={finish} />
          </Col>
        </Row>
      </React.Fragment>
    );
  }
}

const mapStateToProps = state => {
  const { experiment, user } = state;
  const {
    error,
    loading,
    fileLoading,
    name,
    description,
    assignments,
    currentAssigment,
  } = experiment;

  const filePair = getCurrentFilePair(state);
  const diff = filePair ? filePair.diff : null;

  const assignmentsOptions = assignments.map((a, i) => {
    const status = a.answer ? ` (${a.answer})` : '';
    return { value: a.id, name: `${i + 1}${status}` };
  });

  return {
    error,
    loading,
    name,
    description,
    percent: getProgressPercent(state),
    fileLoading,
    diffString: diff,
    currentAssigmentId: currentAssigment ? currentAssigment.id : null,
    assignmentsOptions,
    user,
  };
};

const mapDispatchToProps = dispatch => ({
  load: () => dispatch(load()),
  markSimilar: () => dispatch(markCurrent(ANSWER_SIMILAR)),
  markMaybe: () => dispatch(markCurrent(ANSWER_MAYBE)),
  markDifferent: () => dispatch(markCurrent(ANSWER_DIFFERENT)),
  skip: () => dispatch(markCurrent(ANSWER_SKIP)),
  selectAssigmentId: id =>
    dispatch(
      push(makeUrl('question', { experiment: experimentId, question: id }))
    ),
  finish: () => dispatch(push(makeUrl('finish', { experiment: experimentId }))),
});

export default connect(mapStateToProps, mapDispatchToProps)(Experiment);
