import React from 'react';
import { Button, ButtonToolbar } from 'react-bootstrap';

function Selector({ finish, skip }) {
  return (
    <ButtonToolbar>
      <Button onClick={finish}>Finish</Button>
      <Button bsStyle="info" onClick={skip}>
        Skip
      </Button>
    </ButtonToolbar>
  );
}

export default Selector;
