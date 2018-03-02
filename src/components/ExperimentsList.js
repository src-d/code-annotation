import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Table, Button, Glyphicon } from 'react-bootstrap';
import { makeUrl } from '../state/routes';
import './ExperimentsList.less';

class ExperimentsList extends Component {
  handleOpenEdit(expId) {
    return () => this.props.onOpenEdit(expId);
  }

  buttonText(progress) {
    switch (progress) {
      case 0:
        return 'Begin';
      case 100:
        return 'Finished';
      default:
        return 'Continue';
    }
  }

  render() {
    return (
      <Table responsive className="experiments-list">
        <tbody>
          {this.props.experiments.map(exp => {
            // const started = exp.progress === 0;
            const finished = exp.progress === 100;
            return (
              <tr
                key={exp.id}
                className={`experiments-list__row${
                  finished ? ' finished' : ''
                }`}
              >
                <td className="experiments-list__name">{exp.name}</td>
                {this.props.user.role === 'requester' && (
                  <td className="experiments-list__edit">
                    <Button
                      bsStyle="link"
                      onClick={this.handleOpenEdit(exp.id)}
                    >
                      <Glyphicon glyph="edit" />
                    </Button>
                  </td>
                )}
                <td className="experiments-list__progress">
                  {Math.round(exp.progress)}%
                </td>
                <td className="experiments-list__additional-actions">
                  {/*
                    {!finished && !started ? (
                        <a href="/">mark as finished</a>
                      ) : null}
                  */}
                </td>
                <td className="experiments-list__actions">
                  <Button
                    bsStyle="primary"
                    bsSize="xsmall"
                    disabled={finished}
                    href={makeUrl('experiment', { experiment: exp.id })}
                  >
                    {this.buttonText(exp.progress)}
                  </Button>
                </td>
              </tr>
            );
          })}
        </tbody>
      </Table>
    );
  }
}

const mapStateToProps = state => ({
  user: state.user,
  experiments: state.experiments.list,
});

export default connect(mapStateToProps)(ExperimentsList);
