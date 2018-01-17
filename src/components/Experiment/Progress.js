import React from 'react';
import { ProgressBar } from 'react-bootstrap';
import './Progress.less';

function Progress({ percent, className = '' }) {
  return (
    <div className={`ex-progress ${className}`}>
      <div className="ex-progress__percent">{percent}%</div>
      <div className="ex-progress__bar">
        <ProgressBar striped bsStyle="info" now={percent} />
      </div>
    </div>
  );
}

export default Progress;
