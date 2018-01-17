import React, { Component } from 'react';
import { connect } from 'react-redux';
import { push } from 'redux-little-router';
import { Modal, Button } from 'react-bootstrap';
import PageHeader from '../components/PageHeader';
import Loader from '../components/Loader';
import Progress from '../components/Experiment/Progress';
import Diff from '../components/Experiment/Diff';
import Footer from '../components/Experiment/Footer';
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

const ALWAYS_ALLOW_SKIP_KEY = 'ca/ALWAYS_ALLOW_SKIP';

class Experiment extends Component {
  constructor(props) {
    super(props);
    this.state = {
      showModal: !!window.localStorage.getItem(ALWAYS_ALLOW_SKIP_KEY),
      alwaysAllowSkip: false,
    };
    this.showModal = this.showModal.bind(this);
    this.hideModal = this.hideModal.bind(this);
    this.onAllowSkipChange = this.onAllowSkipChange.bind(this);
  }

  showModal() {
    this.setState({ showModal: true });
  }

  hideModal() {
    this.setState({ showModal: false });
  }

  onAllowSkipChange(e) {
    this.setState({ alwaysAllowSkip: e.target.checked });
    if (e.target.checked) {
      window.localStorage.setItem(ALWAYS_ALLOW_SKIP_KEY, e.target.checked);
    } else {
      window.localStorage.removeItem(ALWAYS_ALLOW_SKIP_KEY);
    }
  }

  render() {
    const { user } = this.props;

    return (
      <div className="ex-page">
        <PageHeader {...user} />
        {this.renderMain()}
        {this.renderModal()}
      </div>
    );
  }

  renderMain() {
    const { error, loading, percent } = this.props;

    if (error) {
      return (
        <div className="ex-page__oops">
          Oops.<br />Something went wrong.
        </div>
      );
    }

    if (loading) {
      return (
        <div className="ex-page__loader">
          <Loader />
        </div>
      );
    }

    return (
      <div className="ex-page__main">
        <div className="ex-page__progress">
          <Progress percent={percent} className="pull-right" />
        </div>
        {this.renderContent()}
      </div>
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
      <div className="ex-page__content">
        <Diff diffString={diffString} className="ex-page__diff" />
        <Footer
          options={assignmentsOptions}
          value={currentAssigmentId}
          select={selectAssigmentId}
          markSimilar={markSimilar}
          markMaybe={markMaybe}
          markDifferent={markDifferent}
          skip={() => (this.state.alwaysAllowSkip ? skip() : this.showModal())}
          finish={finish}
        />
      </div>
    );
  }

  renderModal() {
    const { skip } = this.props;

    return (
      <Modal show={this.state.showModal} onHide={this.hideModal}>
        <Modal.Header closeButton>
          <Modal.Title>
            Are you sure you want to skip this question?
          </Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <p>You can always come back to this question later.</p>
          <label>
            <input
              type="checkbox"
              checked={this.state.alwaysAllowSkip}
              onChange={this.onAllowSkipChange}
            />{' '}
            always allow skiping questions
          </label>
        </Modal.Body>
        <Modal.Footer>
          <Button onClick={this.hideModal}>No</Button>
          <Button
            bsStyle="primary"
            onClick={() => {
              this.hideModal();
              skip();
            }}
          >
            Yes
          </Button>
        </Modal.Footer>
      </Modal>
    );
  }
}

const mapStateToProps = state => {
  const { experiment, user } = state;
  const {
    error,
    loading,
    fileLoading,
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
