import React from 'react';
import { Table, Button } from 'react-bootstrap';
import './ExperimentsList.less';

function buttonText(procent) {
  switch (procent) {
    case 0:
      return 'Begin';
    case 100:
      return 'Finished';
    default:
      return 'Continue';
  }
}

function ExperimentsList({ experiments }) {
  return (
    <Table responsive className="experiments-list">
      <tbody>
        {experiments.map(exp => {
          const started = exp.procent === 0;
          const finished = exp.procent === 100;
          return (
            <tr
              key={exp.id}
              className={`experiments-list__row${finished ? ' finished' : ''}`}
            >
              <td className="experiments-list__name">{exp.name}</td>
              <td className="experiments-list__procent">{exp.procent}%</td>
              <td className="experiments-list__additional-actions">
                {!finished && !started ? (
                  <a href="#">mark as finished</a>
                ) : null}
              </td>
              <td className="experiments-list__actions">
                <Button bsStyle="primary" bsSize="xsmall" disabled={finished}>
                  {buttonText(exp.procent)}
                </Button>
              </td>
            </tr>
          );
        })}
      </tbody>
    </Table>
  );
}

export default ExperimentsList;
