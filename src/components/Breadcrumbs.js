import React from 'react';
import './Breadcrumbs.less';

function Breadcrumbs({ items }) {
  return (
    <ul className="breadcrumbs">
      {items.map((item, i) => (
        <li key={i}>
          <a href={item.link}>{item.name}</a>
        </li>
      ))}
    </ul>
  );
}

export default Breadcrumbs;
