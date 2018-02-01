import React, { Component } from 'react';
import {
  Grid,
  Row,
  Col,
  FormGroup,
  Radio,
  Tabs,
  Tab,
  Table,
} from 'react-bootstrap';
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
    return (
      <Grid fluid className="review-results">
        <Row className="review-results__info">
          <Col xs={4}>
            <FormGroup>
              <Radio
                inline
                value="most"
                checked={this.state.mode === 'most'}
                onChange={this.onModeChange}
              >
                Most similar & disimilar Features
              </Radio>{' '}
              <Radio
                inline
                value="full"
                checked={this.state.mode === 'full'}
                onChange={this.onModeChange}
              >
                Full Output
              </Radio>
            </FormGroup>
          </Col>
          <Col xs={4} className="text-right">
            Users Annotations: <b>24 Similar</b> & <b>3 Disimilar</b>
          </Col>
          <Col xs={4}>
            <div className="review-results__score">
              Similarity Score: 0.9512810301340767
            </div>
          </Col>
        </Row>
        <Row className="review-results__main">
          <Tabs defaultActiveKey={1} id="tabs" className="review-results__tabs">
            <Tab eventKey={1} title="node2vec">
              <div className="review-results__tab">
                <Table
                  striped
                  bordered
                  condensed
                  hover
                  className="review-results__table"
                >
                  <thead>
                    <tr>
                      <th className="review-results__th-number" />
                      <th className="review-results__th-title">
                        Most Similar Features
                      </th>
                      <th>Value One</th>
                      <th>Value Two</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr>
                      <td>1</td>
                      <td>
                        r.QualifiedNameIfStatementQualifiedNameIfStatementQualifiedNameSimpleName
                      </td>
                      <td>4.351873508061938</td>
                      <td>4.351873508061938</td>
                    </tr>
                    <tr>
                      <td>2</td>
                      <td>
                        r.BlockVariableDeclarationStatementSimpleTypeVariableDeclarationStatementBlock
                      </td>
                      <td>2.351873508061938</td>
                      <td>2.351873508061938</td>
                    </tr>
                    <tr>
                      <td>3</td>
                      <td>
                        r.ExpressionStatementBlockIfStatementBlockReturnStatement
                      </td>
                      <td>1.351873508061938</td>
                      <td>1.351873508061938</td>
                    </tr>
                    <tr>
                      <td>4</td>
                      <td>
                        r.QualifiedNameIfStatementQualifiedNameIfStatementQualifiedName
                      </td>
                      <td>0.351873508061938</td>
                      <td>0.351873508061938</td>
                    </tr>
                    <tr>
                      <td>5</td>
                      <td>
                        r.BlockVariableDeclarationStatementSimpleTypeVariableDeclarationStatement
                      </td>
                      <td>4.351873508061938</td>
                      <td>4.351873508061938</td>
                    </tr>
                    <tr>
                      <td>6</td>
                      <td>
                        r.QualifiedNameIfStatementQualifiedNameIfStatementQualifiedNameSimpleName
                      </td>
                      <td>2.351873508061938</td>
                      <td>2.351873508061938</td>
                    </tr>
                    <tr>
                      <td>7</td>
                      <td>
                        r.BlockVariableDeclarationStatementSimpleTypeVariableDeclarationBlock
                      </td>
                      <td>1.351873508061938</td>
                      <td>1.351873508061938</td>
                    </tr>
                    <tr>
                      <td>8</td>
                      <td>
                        r.ExpressionStatementBlockIfStatementBlockReturnStatement
                      </td>
                      <td>0.351873508061938</td>
                      <td>0.351873508061938</td>
                    </tr>
                  </tbody>
                </Table>
                <Table
                  striped
                  bordered
                  condensed
                  hover
                  className="review-results__table"
                >
                  <thead>
                    <tr>
                      <th />
                      <th>Most Disimilar Features</th>
                      <th>Value One</th>
                      <th>Value Two</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr>
                      <td>1</td>
                      <td>
                        r.QualifiedNameIfStatementQualifiedNameIfStatementQualifiedNameSimpleName
                      </td>
                      <td>4.351873508061938</td>
                      <td>4.351873508061938</td>
                    </tr>
                    <tr>
                      <td>2</td>
                      <td>
                        r.BlockVariableDeclarationStatementSimpleTypeVariableDeclarationStatementBlock
                      </td>
                      <td>2.351873508061938</td>
                      <td>2.351873508061938</td>
                    </tr>
                    <tr>
                      <td>3</td>
                      <td>
                        r.ExpressionStatementBlockIfStatementBlockReturnStatement
                      </td>
                      <td>1.351873508061938</td>
                      <td>1.351873508061938</td>
                    </tr>
                    <tr>
                      <td>4</td>
                      <td>
                        r.QualifiedNameIfStatementQualifiedNameIfStatementQualifiedName
                      </td>
                      <td>0.351873508061938</td>
                      <td>0.351873508061938</td>
                    </tr>
                    <tr>
                      <td>5</td>
                      <td>
                        r.BlockVariableDeclarationStatementSimpleTypeVariableDeclarationStatement
                      </td>
                      <td>4.351873508061938</td>
                      <td>4.351873508061938</td>
                    </tr>
                    <tr>
                      <td>6</td>
                      <td>
                        r.QualifiedNameIfStatementQualifiedNameIfStatementQualifiedNameSimpleName
                      </td>
                      <td>2.351873508061938</td>
                      <td>2.351873508061938</td>
                    </tr>
                    <tr>
                      <td>7</td>
                      <td>
                        r.BlockVariableDeclarationStatementSimpleTypeVariableDeclarationBlock
                      </td>
                      <td>1.351873508061938</td>
                      <td>1.351873508061938</td>
                    </tr>
                    <tr>
                      <td>8</td>
                      <td>
                        r.ExpressionStatementBlockIfStatementBlockReturnStatement
                      </td>
                      <td>0.351873508061938</td>
                      <td>0.351873508061938</td>
                    </tr>
                  </tbody>
                </Table>
              </div>
            </Tab>
            <Tab eventKey={2} title="uast2seq">
              Fill it with real data
            </Tab>
            <Tab eventKey={3} title="hash">
              Fill it with real data
            </Tab>
            <Tab eventKey={4} title="other extraction method">
              Fill it with real data
            </Tab>
          </Tabs>
        </Row>
      </Grid>
    );
  }
}

export default Results;
