import React, { Component } from 'react';
import { Helmet } from 'react-helmet';
import PageHeader from '../components/PageHeader';

function Page(WrappedComponent, { className, titleFn, showHeader = false }) {
  class BasePage extends Component {
    constructor(props) {
      super(props);

      this.state = { showHeader };
      this.showHeader = this.showHeader.bind(this);
      this.hideHeader = this.hideHeader.bind(this);
    }

    render() {
      return (
        <div className={className}>
          {titleFn && (
            <Helmet>
              <title>{titleFn(this.props)}</title>
            </Helmet>
          )}
          {this.state.showHeader && <PageHeader />}
          <WrappedComponent
            {...this.props}
            showHeader={this.showHeader}
            hideHeader={this.hideHeader}
          />
        </div>
      );
    }

    showHeader() {
      this.setState({ showHeader: true });
    }

    hideHeader() {
      this.setState({ showHeader: false });
    }
  }

  return BasePage;
}

export default Page;
