import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Grid, Row, Col } from 'react-bootstrap';
import SplitPane from 'react-split-pane';
import PageHeader from '../components/PageHeader';
import Breadcrumbs from '../components/Breadcrumbs';
import Selector from '../components/Experiment/Selector';
import Diff from '../components/Experiment/Diff';
import Results from '../components/Review/Results';
import diffString from '../api/respMock';
import './Review.less';

class Review extends Component {
  render() {
    const { user } = this.props;

    return (
      <div className="review-page">
        <PageHeader {...user} />
        <Grid fluid className="review-page__grid">
          <Row className="review-page__header">
            <Col xs={8} className="review-page__breadcrumbs">
              <Breadcrumbs
                items={[
                  { name: 'some', link: '#' },
                  { name: 'bread', link: '#' },
                  { name: 'review', link: '#' },
                ]}
              />
            </Col>
            <Col xs={4} className="review-page__controls">
              <Selector
                title="Previous"
                options={[{ value: 1, name: '(1) (Yes)' }]}
              />
            </Col>
          </Row>
          <Row className="review-page__main-row">
            <Col xs={12} className="review-page__main-col">
              {/* we need this wrapper because SplitPane hard coded width 100% */}
              <div className="review-page__split">
                <SplitPane
                  split="horizontal"
                  minSize={50}
                  defaultSize={100}
                  style={{ width: 'auto' }}
                >
                  <Diff diffString={diffString} className="review-page__diff" />
                  <Results className="results" />
                </SplitPane>
              </div>
            </Col>
          </Row>
        </Grid>
      </div>
    );
  }
}

const mapStateToProps = state => ({
  user: state.user,
});

export default connect(mapStateToProps)(Review);
