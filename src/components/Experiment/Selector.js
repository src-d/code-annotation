import React, { Component } from 'react';
import './Selector.less';

class Selector extends Component {
  constructor(props) {
    super(props);

    this.titleEl = null;
    this.state = {
      titleWidth: 0,
    };
  }

  componentDidMount() {
    this.setState({ titleWidth: this.titleEl.offsetWidth });
  }

  componentDidUpdate(prevProps) {
    if (this.props.title === prevProps.title) {
      return;
    }
    this.setState({ titleWidth: this.titleEl.offsetWidth });
  }

  render() {
    const { title, value, onChange, options } = this.props;

    return (
      <div className="selector">
        <span
          className="selector__title"
          ref={e => {
            this.titleEl = e;
          }}
        >
          {title}
        </span>
        <select
          value={value}
          onChange={onChange}
          ref={e => {
            this.select = e;
          }}
          style={{ paddingLeft: this.state.titleWidth }}
        >
          {options.map(opt => (
            <option value={opt.value} key={opt.value}>
              {opt.name}
            </option>
          ))}
        </select>
      </div>
    );
  }
}

export default Selector;
