import React, { Component } from 'react';
import { connect } from 'react-redux';

import {
  Button,
  Modal,
  FormGroup,
  FormControl,
  Glyphicon,
  Row,
  Col,
} from 'react-bootstrap';

import {
  create as experimentCreate,
  update as experimentUpdate,
  uploadFilePairs,
  UPLOAD_RES_SUCCESS,
  UPLOAD_RES_FAILURE,
} from '../../state/experiments';

import './ExperimentEditModal.less';

class ExperimentEditModal extends Component {
  constructor(props) {
    super(props);

    this.state = this.buildState(this.props);

    this.handleNameChange = this.handleNameChange.bind(this);
    this.handleDescriptionChange = this.handleDescriptionChange.bind(this);
    this.handleSave = this.handleSave.bind(this);
    this.handleUpload = this.handleUpload.bind(this);
  }

  buildState(props) {
    const exp = props.editModalExperiment;

    if (props.action === 'edit' && exp !== undefined) {
      return { nameVal: exp.name, descriptionVal: exp.description };
    }

    return { nameVal: '', descriptionVal: '' };
  }

  handleNameChange(e) {
    this.setState({ nameVal: e.target.value });
  }

  handleDescriptionChange(e) {
    this.setState({ descriptionVal: e.target.value });
  }

  handleSave() {
    const create = this.props.action === 'create';

    if (create) {
      return this.props
        .experimentCreate(this.state.nameVal, this.state.descriptionVal)
        .then(expId => {
          this.props.onOpenEdit(expId);
        })
        .catch(() => {});
    }

    return this.props.experimentUpdate(
      this.props.editModalExperiment.id,
      this.state.nameVal,
      this.state.descriptionVal
    );
  }

  handleUpload(e) {
    this.props.uploadFilePairs(
      this.props.editModalExperiment.id,
      e.target.files[0]
    );
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.editModalExperiment !== this.props.editModalExperiment) {
      this.setState(this.buildState(nextProps));
    }
  }

  render() {
    const {
      createInProgress,
      updateInProgress,
      uploadInProgress,
      uploadResult,
    } = this.props.experiments;

    const create = this.props.action === 'create';

    const title = create ? 'Create Experiment' : 'Edit Experiment';
    const btn = create ? 'Create' : 'Save';
    const btnDisabled =
      createInProgress || updateInProgress || this.state.nameVal.trim() === '';

    const uploadDisabled = create || uploadInProgress;

    let uploadIcon;
    let uploadClass;
    if (uploadInProgress) {
      uploadIcon = 'floppy-open';
      uploadClass = 'upload-floppy upload-floppy-progress';
    } else if (uploadResult === UPLOAD_RES_SUCCESS) {
      uploadIcon = 'floppy-saved';
      uploadClass = 'upload-floppy upload-floppy-success';
    } else if (uploadResult === UPLOAD_RES_FAILURE) {
      uploadIcon = 'floppy-remove';
      uploadClass = 'upload-floppy upload-floppy-failure';
    } else {
      uploadIcon = 'floppy-disk';
      uploadClass = 'upload-floppy';
    }

    return (
      <Modal show={this.props.show} onHide={this.props.onHide}>
        <Modal.Header closeButton>
          <Modal.Title>{title}</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <FormGroup controlId="nameInput">
            <FormControl
              type="text"
              value={this.state.nameVal}
              placeholder="Experiment name"
              onChange={this.handleNameChange}
            />
          </FormGroup>
          <FormGroup controlId="descInput">
            <FormControl
              type="text"
              value={this.state.descriptionVal}
              placeholder="Experiment description"
              onChange={this.handleDescriptionChange}
            />
          </FormGroup>
          <hr />
          <Row>
            <Col xs={7}>
              <p>
                <Glyphicon className={uploadClass} glyph={uploadIcon} />
                <b> Experiment Dataset</b> SQLite DB file
              </p>
            </Col>
            <Col xs={5} className="text-right">
              <FormGroup controlId="uploadTest">
                <FormControl
                  type="file"
                  onChange={this.handleUpload}
                  disabled={uploadDisabled}
                />
              </FormGroup>
            </Col>
          </Row>
        </Modal.Body>
        <Modal.Footer>
          <Button
            bsStyle="primary"
            disabled={btnDisabled}
            onClick={this.handleSave}
          >
            {btn}
          </Button>
        </Modal.Footer>
      </Modal>
    );
  }
}

const mapStateToProps = state => ({
  experiments: state.experiments,
});

export default connect(mapStateToProps, {
  experimentCreate,
  experimentUpdate,
  uploadFilePairs,
})(ExperimentEditModal);
