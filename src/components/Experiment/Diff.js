import React, { PureComponent } from 'react';
import { Diff2Html } from 'diff2html';
import 'diff2html/dist/diff2html.css';
import './Diff.less';

class Diff extends PureComponent {
  render() {
    const { diffString, leftLoc, rightLoc, className } = this.props;
    const diffHTML = Diff2Html.getPrettyHtml(diffString, {
      inputFormat: 'diff',
      outputFormat: 'side-by-side',
      showFiles: false,
      matching: 'none',
      matchWordsThreshold: 0.25,
      matchingMaxComparisons: 2500,
    });
    return (
      <div className={`diff ${className}`}>
        <div className="diff__locs">
          <div className="diff__loc left">{leftLoc} lines of code</div>
          <div className="diff__loc right">{rightLoc} lines of code</div>
        </div>
        <div
          className="diff__content"
          dangerouslySetInnerHTML={{ __html: diffHTML }}
        />
      </div>
    );
  }
}

export default Diff;
