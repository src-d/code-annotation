import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Grid, Row, Col, Button, Glyphicon } from 'react-bootstrap';
import Page from './Page';
import Loader from '../components/Loader';
import ExperimentsList from '../components/ExperimentsList';
import ExperimentsEditModal from '../components/Experiments/ExperimentEditModal';
import {
  load as experimentsLoad,
  uploadResultReset,
} from '../state/experiments';
import './Experiments.less';

class Experiments extends Component {
  constructor(props) {
    super(props);

    this.handleOpenNew = this.handleOpenNew.bind(this);
    this.handleOpenEdit = this.handleOpenEdit.bind(this);
    this.handleCloseModal = this.handleCloseModal.bind(this);

    this.state = {
      editModalShow: false,
      editModalAction: 'create',
      editExperimentId: undefined,
    };
  }

  componentDidMount() {
    this.props.load();
  }

  handleOpenNew() {
    this.setState({
      editModalShow: true,
      editModalAction: 'create',
    });
  }

  handleOpenEdit(experimentId) {
    this.props.uploadResultReset();

    this.setState({
      editModalShow: true,
      editModalAction: 'edit',
      editExperimentId: experimentId,
    });
  }

  handleCloseModal() {
    this.setState({ editModalShow: false });
  }

  render() {
    return <Grid>{this.renderContent()}</Grid>;
  }

  renderContent() {
    const { user, experiments } = this.props;
    const { error, loading, createInProgress } = experiments;

    if (error) {
      return (
        <Row className="ex-page__oops">
          <Col xs={12}>
            Oops.<br />Something went wrong.
          </Col>
        </Row>
      );
    }

    if (loading) {
      return (
        <Row className="ex-page__loader">
          <Col xs={12}>
            <Loader />
          </Col>
        </Row>
      );
    }

    let addRow = '';
    if (user.role === 'requester') {
      addRow = (
        <Row>
          <Col xs={8} xsOffset={2} className="text-right">
            <Button
              bsStyle="link"
              onClick={this.handleOpenNew}
              disabled={createInProgress}
              className="ex-page__plus-btn"
            >
              <Glyphicon glyph="plus-sign" />
            </Button>
          </Col>
        </Row>
      );
    }

    let editModalExperiment;

    if (this.state.editModalAction === 'edit') {
      editModalExperiment = experiments.list.find(
        e => e.id === this.state.editExperimentId
      );
    }

    return (
      <React.Fragment>
        <Row>
          <Col xs={12}>
            <h1 className="text-center">
              You look great today, {user.username}!
            </h1>
          </Col>
        </Row>
        <Row>
          <Col xs={12} style={{ paddingBottom: '40px' }}>
            <p className="text-center">
              Here&#39;s your experiments dashboard, go ahead and make our day.
            </p>
          </Col>
        </Row>
        {addRow}
        <Row>
          <Col xs={8} xsOffset={2}>
            <ExperimentsList onOpenEdit={this.handleOpenEdit} />
          </Col>
        </Row>
        <ExperimentsEditModal
          action={this.state.editModalAction}
          editModalExperiment={editModalExperiment}
          show={this.state.editModalShow}
          onHide={this.handleCloseModal}
          onOpenEdit={this.handleOpenEdit}
        />
      </React.Fragment>
    );
  }
}

const mapStateToProps = state => ({
  user: state.user,
  experiments: state.experiments,
});

export default connect(mapStateToProps, {
  load: experimentsLoad,
  uploadResultReset,
})(
  Page(Experiments, {
    className: 'experiments-page',
    titleFn: () => 'Experiments',
    showHeader: true,
  })
);
