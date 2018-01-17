import React from 'react';
import { Button } from 'react-bootstrap';
import './Footer.less';

function Footer({
  options,
  value,
  select,
  markSimilar,
  markMaybe,
  markDifferent,
  skip,
  finish,
}) {
  return (
    <div className="ex-footer">
      <div className="container">
        <div className="ex-footer__pager">
          Previous:{' '}
          <select value={value} onChange={e => select(e.target.value)}>
            {options.map(opt => (
              <option value={opt.value} key={opt.value}>
                {opt.name}
              </option>
            ))}
          </select>
        </div>
        <div className="ex-footer__actions">
          <Button bsSize="large" bsStyle="success" onClick={markSimilar}>
            Identical
          </Button>
          <Button bsSize="large" onClick={markMaybe}>
            Similar
          </Button>
          <Button bsSize="large" bsStyle="danger" onClick={markDifferent}>
            Different
          </Button>
        </div>
        <div className="ex-footer__buttons">
          <Button onClick={finish}>Finish</Button>
          <Button bsStyle="info" onClick={skip}>
            Skip
          </Button>
        </div>
      </div>
    </div>
  );
}

export default Footer;
