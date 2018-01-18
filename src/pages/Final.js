import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Table } from 'react-bootstrap';
import PageHeader from '../components/PageHeader';
import HumanDuration from '../components/HumanDuration';
import {
  getSimilarCount,
  getDifferentCount,
  getProgressPercent,
  getSkipCount,
  getAverageTime,
  getOverallTime,
} from '../state/experiment';
import './Final.less';

class Final extends Component {
  render() {
    const {
      user,
      expName,
      procent,
      similarCount,
      differentCount,
      skipCount,
      averageTime,
      overallTime,
    } = this.props;

    return (
      <div className="final-page">
        <PageHeader {...user} />
        <h1>You did it great, {user.username}!</h1>
        <p>Here&#39;s an overview of {expName} experiment</p>
        <Table bordered className="final-page__table">
          <tbody>
            <tr>
              <td>Completion</td>
              <td>{procent}%</td>
            </tr>
            <tr>
              <td>Similar Annotations</td>
              <td>{similarCount}</td>
            </tr>
            <tr>
              <td>Different Annotations</td>
              <td>{differentCount}</td>
            </tr>
            <tr>
              <td>Skiped Annotations</td>
              <td>{skipCount}</td>
            </tr>
            <tr>
              <td>Average Time per Annotation</td>
              <td>
                <HumanDuration value={averageTime} />
              </td>
            </tr>
            <tr>
              <td>Experiment Overall Time</td>
              <td>
                <HumanDuration value={overallTime} />
              </td>
            </tr>
          </tbody>
        </Table>
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
    similarCount: getSimilarCount(state),
    differentCount: getDifferentCount(state),
    skipCount: getSkipCount(state),
    averageTime: getAverageTime(state),
    overallTime: getOverallTime(state),
  };
};

export default connect(mapStateToProps)(Final);
