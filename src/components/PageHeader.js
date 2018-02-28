import React from 'react';
import { Navbar, Nav, NavItem, NavDropdown, MenuItem } from 'react-bootstrap';
import { connect } from 'react-redux';
import { push } from 'redux-little-router';
import { logOut } from '../state/user';
import { makeUrl } from '../state/routes';
import './PageHeader.less';

function PageHeader({
  username,
  role,
  avatarUrl,
  onReviewClick,
  onExportClick,
  onLogoutClick,
}) {
  const menuItems = [
    { name: 'Review', onClick: onReviewClick, role: 'requester' },
    { name: 'Export', onClick: onExportClick, role: 'requester' },
    { role: 'requester' },
    { name: 'Logout', onClick: onLogoutClick },
  ];

  const items = menuItems.reduce((acc, item, i) => {
    if (item.role && item.role !== role) {
      return acc;
    }
    if (!item.name) {
      return [...acc, <MenuItem divider key={i} />];
    }
    return [
      ...acc,
      <MenuItem eventKey={i} onClick={item.onClick} key={i}>
        {item.name}
      </MenuItem>,
    ];
  }, []);

  return (
    <Navbar fluid className="header">
      <Navbar.Header>
        <Navbar.Brand>
          <h1>
            <a href="/">
              <img
                alt="source{d} code annotation"
                src="/sourced-code-annotation-logo.svg"
              />
              code annotation
            </a>
          </h1>
        </Navbar.Brand>
      </Navbar.Header>
      <Nav pullRight>
        <NavDropdown eventKey={1} title="dashboard" id="nav-dropdown">
          {items}
        </NavDropdown>
        <NavItem eventKey={2}>
          <img
            src={avatarUrl}
            alt={username}
            style={{ width: '24px', height: '24px', margin: '-3px 0' }}
          />
        </NavItem>
      </Nav>
    </Navbar>
  );
}

export default connect(state => state.user, {
  onReviewClick: () => push(makeUrl('review', { experiment: 1 })),
  onExportClick: () => push(makeUrl('export')),
  onLogoutClick: logOut,
})(PageHeader);
