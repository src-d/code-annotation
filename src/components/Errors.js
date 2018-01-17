import React, { PureComponent } from 'react';
import { Alert } from 'react-bootstrap';
import { connect } from 'react-redux';
import { remove } from '../state/errors';
import './Errors.less';

class Errors extends PureComponent {
  render() {
    const { errors, onDismiss } = this.props;
    return (
      <div className="errors">
        {errors.map((e, i) => (
          <Alert bsStyle="danger" key={i} onDismiss={() => onDismiss(i)}>
            {e}
          </Alert>
        ))}
      </div>
    );
  }
}

const mapStateToProps = state => ({ errors: state.errors });

export default connect(mapStateToProps, { onDismiss: remove })(Errors);
