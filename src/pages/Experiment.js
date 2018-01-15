import React, { Component } from 'react';
import { connect } from 'react-redux';
import PageHeader from '../components/PageHeader';
import Loader from '../components/Loader';
import Progress from '../components/Experiment/Progress';
import Diff from '../components/Experiment/Diff';
import Footer from '../components/Experiment/Footer';
import { load, selectAssigment, markCurrent } from '../state/experiment';
import './Experiment.less';

class Experiment extends Component {
  componentDidMount() {
    this.props.load();
  }

  render() {
    const {
      loading,
      stat,
      fileLoading,
      diffString,
      assignmentsOptions,
      currentAssigmentId,
      user,
      selectAssigmentId,
      markSimilar,
      markMaybe,
      markDifferent,
      skip,
    } = this.props;
    return (
      <div className="ex-page">
        <PageHeader {...user} />
        {loading ? (
          <div className="ex-page__loader">
            <Loader />
          </div>
        ) : (
          <div className="ex-page__main">
            <div className="ex-page__progress">
              <Progress
                percent={stat.percent}
                similar={stat.similar}
                different={stat.different}
                className="pull-right"
              />
            </div>
            {fileLoading ? (
              <div className="ex-page__loader">
                <Loader />
              </div>
            ) : (
              <div className="ex-page__content">
                <Diff diffString={diffString} className="ex-page__diff" />
                <Footer
                  options={assignmentsOptions}
                  value={currentAssigmentId}
                  select={selectAssigmentId}
                  markSimilar={markSimilar}
                  markMaybe={markMaybe}
                  markDifferent={markDifferent}
                  skip={skip}
                />
              </div>
            )}
          </div>
        )}
      </div>
    );
  }
}

const mapStateToProps = state => {
  const { experiment, user } = state;
  const { loading, assignments, currentAssigment } = experiment;

  const similar = assignments.filter(a => a.answer === 'Yes').length;
  const different = assignments.filter(a => a.answer === 'No').length;
  const stat = {
    similar,
    different,
    percent: 100 / assignments.length * (similar + different),
  };

  let diff = null;
  if (!loading && currentAssigment) {
    const filePair = experiment.filePairs[currentAssigment.pairId];
    ({ diff } = filePair);
  }

  const assignmentsOptions = assignments.map((a, i) => {
    const status = a.answer ? ` (${a.answer})` : '';
    return { value: a.id, name: `${i + 1}${status}` };
  });

  return {
    loading,
    stat,
    fileLoading: experiment.fileLoading,
    diffString: diff,
    currentAssigmentId: currentAssigment ? currentAssigment.id : null,
    assignmentsOptions,
    user,
  };
};

const mapDispatchToProps = dispatch => ({
  load: () => dispatch(load()),
  markSimilar: () => dispatch(markCurrent('Yes')),
  markMaybe: () => dispatch(markCurrent('Maybe')),
  markDifferent: () => dispatch(markCurrent('No')),
  skip: () => dispatch(markCurrent('Skip')),
  selectAssigmentId: id => dispatch(selectAssigment(id)),
});

export default connect(mapStateToProps, mapDispatchToProps)(Experiment);
