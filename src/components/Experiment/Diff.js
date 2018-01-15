import React, { Component } from 'react';
import { Diff2Html } from 'diff2html';
import 'diff2html/dist/diff2html.css';
import './Diff.less';

// TODO add componentWillReceive props
// we need to minimize rerender of this component, because it's very heavy one
class Diff extends Component {
  render() {
    const { diffString, className } = this.props;
    const diffHTML = Diff2Html.getPrettyHtml(diffString, {
      inputFormat: 'diff',
      outputFormat: 'side-by-side',
      showFiles: false,
      matching: 'none',
      matchWordsThreshold: 0.25,
      matchingMaxComparisons: 2500,
    });
    return (
      <div
        className={className}
        dangerouslySetInnerHTML={{ __html: diffHTML }}
      />
    );
  }
}

export default Diff;
