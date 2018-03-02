import React, { Component } from 'react';
import { Grid, Row, Col, FormGroup, Radio, Tabs, Tab } from 'react-bootstrap';
import FeaturesTable from './FeaturesTable';
import './Results.less';

class Results extends Component {
  constructor(props) {
    super(props);

    this.state = {
      mode: 'most',
    };

    this.onModeChange = this.onModeChange.bind(this);
  }

  onModeChange(e) {
    this.setState({ mode: e.target.value });
  }

  render() {
    const {
      score,
      annotations,
      mostSimilarFeatures,
      leastSimilarFeatures,
      features,
    } = this.props;

    return (
      <Grid fluid className="review-results">
        <Row className="review-results__info">
          <Col xs={12} className="review-results__info-col">
            <FormGroup>
              <Radio
                inline
                value="most"
                checked={this.state.mode === 'most'}
                onChange={this.onModeChange}
                className={this.state.mode === 'most' ? '_checked' : ''}
              >
                Most similar & dissimilar Features
              </Radio>{' '}
              <Radio
                inline
                value="full"
                checked={this.state.mode === 'full'}
                onChange={this.onModeChange}
                className={this.state.mode === 'full' ? '_checked' : ''}
              >
                Full Output
              </Radio>
            </FormGroup>
            <div className="review-results__stats">
              Users Annotations: <b>{annotations.yes || 0} Similar</b> &{' '}
              <b>{annotations.no || 0} Dissimilar</b>
            </div>
            <div className="review-results__score">
              Similarity Score: {score}
            </div>
          </Col>
        </Row>
        <Row className="review-results__main">
          <Tabs defaultActiveKey={1} id="tabs" className="review-results__tabs">
            <Tab eventKey={1} title="All features">
              <div className="review-results__tab">
                {this.state.mode === 'most' ? (
                  <React.Fragment>
                    <FeaturesTable
                      title="Most Similar Features"
                      features={mostSimilarFeatures}
                    />
                    <FeaturesTable
                      title="Most Dissimilar Features"
                      features={leastSimilarFeatures}
                      bsStyle="danger"
                    />
                  </React.Fragment>
                ) : (
                  <FeaturesTable title="Features" features={features} />
                )}
              </div>
            </Tab>
          </Tabs>
        </Row>
      </Grid>
    );
  }
}

export default Results;
