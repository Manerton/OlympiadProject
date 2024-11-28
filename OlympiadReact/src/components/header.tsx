import React, { useState, useEffect } from "react";
import { Navbar, Nav, Dropdown, Container } from "react-bootstrap";
import { LinkContainer } from "react-router-bootstrap";
import { PersonCircle, BoxArrowInRight, BoxArrowInLeft } from "react-bootstrap-icons";

function Header() {
  // Simulating authentication state
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  // Check authentication (e.g., from localStorage or API)
  useEffect(() => {
    const userToken = localStorage.getItem("userToken"); // Example: check for a token
    setIsAuthenticated(!!userToken);
  }, []);

  // Logout function
  const handleLogout = () => {
    localStorage.removeItem("userToken"); // Clear token or session
    setIsAuthenticated(false);
  };

  return (
    <Navbar bg="light" expand="lg">
      <Container fluid>
        <Navbar.Brand href="/">ВСОШ</Navbar.Brand>
        <Navbar.Toggle aria-controls="basic-navbar-nav" />
        <Navbar.Collapse id="basic-navbar-nav" className="navbar-expand-lg">
          <Nav className="me-auto">
            <LinkContainer to="/events">
              <Nav.Link>События</Nav.Link>
            </LinkContainer>
            {/* тут role == 1,2 */}
            <LinkContainer to="/applications">
              <Nav.Link>Статус заявки</Nav.Link>
            </LinkContainer>
            {/* тут role == 3 */}
            <LinkContainer to="/attendance">
              <Nav.Link>Отметить присутствие</Nav.Link>
            </LinkContainer>
          </Nav>
          <Nav className="ms-auto">
            {isAuthenticated ? (
              <Dropdown align="end">
                <Dropdown.Toggle
                  variant="light"
                  id="dropdown-basic"
                  className="d-flex align-items-center border-0 bg-transparent"
                >
                  <PersonCircle size={24} className="me-2" />
                  <span>Имя пользователя</span>
                </Dropdown.Toggle>

                <Dropdown.Menu>
                  <Dropdown.Item href="/profile">
                    <PersonCircle className="me-2" /> Мой профиль
                  </Dropdown.Item>
                  <Dropdown.Item onClick={handleLogout}>
                    <BoxArrowInLeft className="me-2" /> Выйти
                  </Dropdown.Item>
                </Dropdown.Menu>
              </Dropdown>
            ) : (
              <Nav.Link href="/login" className="d-flex align-items-center">
                <BoxArrowInRight className="me-2" />
                Войти
              </Nav.Link>
            )}
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
}

export default Header;
