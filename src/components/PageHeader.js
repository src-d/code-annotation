import React from 'react';
import { Navbar, Nav, NavItem, NavDropdown, MenuItem } from 'react-bootstrap';
import { connect } from 'react-redux';
import { logOut } from '../state/user';
import './PageHeader.less';

function PageHeader({ username, avatarUrl, onLogoutClick }) {
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
          <MenuItem eventKey={1.1} onClick={onLogoutClick}>
            Logout
          </MenuItem>
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

export default connect(state => state.user, { onLogoutClick: logOut })(
  PageHeader
);
