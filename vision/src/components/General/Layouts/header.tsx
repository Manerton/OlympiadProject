import React from "react";
import { Navbar, Nav, Dropdown, Container } from "react-bootstrap";
import { LinkContainer } from "react-router-bootstrap";
import { Link } from "react-router-dom";
import { PersonCircle, BoxArrowInRight } from "react-bootstrap-icons";
import ThemeToggleButton from "../../Helpers/ThemeToggleButton";
import { useAuth } from "../../Helpers/AuthContext";

function Header() {
  const { user, logout } = useAuth();
   console.log(user);
  // Проверка на авторизацию
  const isAuthenticated = user;
 

  return (
    <Navbar expand="lg" sticky="top" className="bg-body-tertiary border border-1">
      <Container fluid="lg" className="mx-auto">
        <Navbar.Brand as={Link} to="/">ВСОШ</Navbar.Brand>
        <Navbar.Toggle aria-controls="basic-navbar-nav" />
        <Navbar.Collapse id="basic-navbar-nav" className="navbar-expand-lg">
          <Nav className="me-auto">
            <LinkContainer to="/">
              <Nav.Link>Главная</Nav.Link>
            </LinkContainer>

            <LinkContainer to="/RegionalStages">
              <Nav.Link>Олимпиады</Nav.Link>
            </LinkContainer>

            {isAuthenticated && (
              <LinkContainer to={`/applications/user/${user?.id}`}>
                <Nav.Link>Статус заявки</Nav.Link>
              </LinkContainer>
            )}

            {/* Админские ссылки, показываются только для админов */}
            {user?.role === 1 && (
              <LinkContainer to="/AdminPanel">
                <Nav.Link>Панель Администрирования</Nav.Link>
              </LinkContainer>
            )}
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
                  <span className="ms-2 text-body">Email-{user?.Email}</span>
                  <span className="ms-2 text-body">id-{user?.id}</span>
                  <span className="ms-2 text-body">role-{user?.role}</span>
                </Dropdown.Toggle>

                <Dropdown.Menu>
                  <Dropdown.Item as={Link} to="/profile">
                    <PersonCircle className="me-2" /> Мой профиль
                  </Dropdown.Item>
                  <Dropdown.Item onClick={logout}>
                    <BoxArrowInRight className="me-2" /> Выйти
                  </Dropdown.Item>
                </Dropdown.Menu>
              </Dropdown>
            ) : (
              <Nav.Link as={Link} to="/auth" className="d-flex align-items-center">
                <BoxArrowInRight className="me-2" />
                Войти
              </Nav.Link>
            )}
          </Nav>
          <ThemeToggleButton />
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
}

export default Header;
