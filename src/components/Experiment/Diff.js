import React, { PureComponent } from 'react';
import { Diff2Html } from 'diff2html';
import 'diff2html/dist/diff2html.css';
import './Diff.less';

class Diff extends PureComponent {
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
