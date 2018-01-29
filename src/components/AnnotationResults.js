import React from 'react';
import { Table } from 'react-bootstrap';
import HumanDuration from '../components/HumanDuration';

function AnnotationResults({
  procent,
  identicalCount,
  similarCount,
  differentCount,
  skipCount,
  averageTime,
  overallTime,
}) {
  return (
    <Table bordered responsive>
      <tbody>
        <tr>
          <td>Completion</td>
          <td>{procent}%</td>
        </tr>
        <tr>
          <td>Identical Annotations</td>
          <td>{identicalCount}</td>
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
  );
}

export default AnnotationResults;
