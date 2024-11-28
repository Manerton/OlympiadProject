import { useState } from "react";
import { Container, Nav, Navbar } from "react-bootstrap"
import { LinkContainer } from 'react-router-bootstrap'


function Header() {
  const [visible, setVisible] = useState(true)
  return (
      <Navbar bg="light" expand="lg">
      <Container fluid>
        <Navbar.Brand href="/">Olympiad Vision</Navbar.Brand>
        <Navbar.Toggle aria-controls="basic-navbar-nav" onClick={() => setVisible(!visible)} />
        <Navbar.Collapse id="basic-navbar-nav" className="navbar-expand-lg">
          <Nav className="me-auto">
            <LinkContainer to="/events">
              <Nav.Link>События</Nav.Link>
            </LinkContainer>
            <LinkContainer to="/jury-assignments">
              <Nav.Link>Назначения жюри</Nav.Link>
            </LinkContainer>
          </Nav>
        </Navbar.Collapse>
      </Container>

      </Navbar>
  )
}

export default Header