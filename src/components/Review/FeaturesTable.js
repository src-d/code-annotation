import React from 'react';
import { Table } from 'react-bootstrap';
import './FeaturesTable.less';

function FeaturesTable({ title, features }) {
  return (
    <Table striped bordered condensed hover className="features-table">
      <thead>
        <tr>
          <th className="features-table__th-number" />
          <th className="features-table__th-title">{title}</th>
          <th>Value One</th>
          <th>Value Two</th>
        </tr>
      </thead>
      <tbody>
        {features.map((f, i) => (
          <tr key={i}>
            <td>{i + 1}</td>
            <td className="features-table__cell-name" title={f.name}>
              {f.name}
            </td>
            <td>{f.weightA}</td>
            <td>{f.weightB}</td>
          </tr>
        ))}
      </tbody>
    </Table>
  );
}

export default FeaturesTable;
