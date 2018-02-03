import React from 'react';
import { Button, ButtonToolbar } from 'react-bootstrap';

function Selector({ markSimilar, markMaybe, markDifferent }) {
  return (
    <ButtonToolbar>
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
  );
}

export default Selector;
