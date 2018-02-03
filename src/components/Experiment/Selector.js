import React from 'react';

function Selector({ options, value, select }) {
  return (
    <div>
      Previous:{' '}
      <select value={value} onChange={e => select(e.target.value)}>
        {options.map(opt => (
          <option value={opt.value} key={opt.value}>
            {opt.name}
          </option>
        ))}
      </select>
    </div>
  );
}

export default Selector;
