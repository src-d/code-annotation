import React from 'react';
import { Grid, Row, Col, Button, ButtonToolbar } from 'react-bootstrap';
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
    <Grid fluid>
      <Row>
        <Col xs={3}>
          Previous:{' '}
          <select value={value} onChange={e => select(e.target.value)}>
            {options.map(opt => (
              <option value={opt.value} key={opt.value}>
                {opt.name}
              </option>
            ))}
          </select>
        </Col>
        <Col xs={6} className="text-center">
          <ButtonToolbar className="ex-footer__main-buttons">
            <Button bsSize="large" bsStyle="success" onClick={markSimilar}>
              Identical
            </Button>
            <Button bsSize="large" onClick={markMaybe}>
              Similar
            </Button>
            <Button bsSize="large" bsStyle="danger" onClick={markDifferent}>
              Different
            </Button>
          </ButtonToolbar>
        </Col>
        <Col xs={3} className="text-right">
          <ButtonToolbar className="ex-footer__additional-buttons">
            <Button onClick={finish}>Finish</Button>
            <Button bsStyle="info" onClick={skip}>
              Skip
            </Button>
          </ButtonToolbar>
        </Col>
      </Row>
    </Grid>
  );
}

export default Footer;
