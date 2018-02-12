import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Grid, Row, Col, Button } from 'react-bootstrap';
import PageHeader from '../components/PageHeader';
import Loader from '../components/Loader';
import { add as addErrors } from '../state/errors';
import api from '../api';

class Export extends Component {
  constructor(props) {
    super(props);

    this.state = {
      loading: true,
      files: [],
    };

    this.loadFiles = this.loadFiles.bind(this);
    this.onCreateClick = this.onCreateClick.bind(this);
  }

  componentDidMount() {
    this.loadFiles().catch(err => {
      this.props.addErrors(err);
    });
  }

  loadFiles() {
    this.setState({ loading: true });

    return api
      .exportList()
      .then(files => {
        this.setState({ files: files || [], loading: false });
      })
      .catch(err => {
        this.setState({ loading: false });
        throw err;
      });
  }

  onCreateClick(e) {
    e.preventDefault();

    return api
      .exportCreate()
      .then(this.loadFiles)
      .catch(err => {
        this.props.addErrors(err);
      });
  }

  render() {
    return (
      <div className="export-page">
        <PageHeader {...this.props.user} />
        <Grid>
          <Row style={{ paddingTop: '20px', paddingBottom: '20px' }}>
            <Col xs={12}>
              <Button
                bsStyle="primary"
                onClick={this.onCreateClick}
                disabled={this.state.loading}
              >
                Create a new SQLite export
              </Button>
            </Col>
          </Row>
          <Row>
            <Col xs={12}>{this.filesList()}</Col>
          </Row>
        </Grid>
      </div>
    );
  }

  filesList() {
    if (this.state.loading) {
      return <Loader />;
    }

    if (!this.state.files.length) {
      return <p>No files to download</p>;
    }

    return (
      <ul>
        {this.state.files.map((f, i) => (
          <li key={i}>
            <a
              href="#"
              onClick={e => {
                e.preventDefault();
                api.exportDownload(f);
              }}
            >
              {f}
            </a>
          </li>
        ))}
      </ul>
    );
  }
}

const mapStateToProps = state => ({
  user: state.user,
});

export default connect(mapStateToProps, { addErrors })(Export);
