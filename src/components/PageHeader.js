import React from 'react';
import { Navbar, Nav, NavItem } from 'react-bootstrap';
import './PageHeader.less';

function PageHeader({ username, avatarUrl }) {
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
        <NavItem eventKey={1} href="/">
          dashboard
        </NavItem>
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

export default PageHeader;
