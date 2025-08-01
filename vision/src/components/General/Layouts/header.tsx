import React from "react";
import { Navbar, Nav, Dropdown, Container } from "react-bootstrap";
import { LinkContainer } from "react-router-bootstrap";
import { PersonCircle, BoxArrowInRight } from "react-bootstrap-icons";

function Header() {
  // Заглушки: в будущем подключите контекст пользователя
  const isAuthenticated = true; // временная проверка
  const userName = "Имя";
  const userId = "123";
  const userRole = "роль";

  return (
    <Navbar bg="light" expand="lg" sticky="top">
      <Container fluid>
        <Navbar.Brand href="/">ВСОШ</Navbar.Brand>
        <Navbar.Toggle aria-controls="basic-navbar-nav" />
        <Navbar.Collapse id="basic-navbar-nav" className="navbar-expand-lg">
          <Nav className="me-auto">

            <LinkContainer to="/">
              <Nav.Link>Главная</Nav.Link>
            </LinkContainer>

            <LinkContainer to="/events">
              <Nav.Link>Олимпиады</Nav.Link>
            </LinkContainer>

            <LinkContainer to={`/applications/user/${userId}`}>
              <Nav.Link>Статус заявки</Nav.Link>
            </LinkContainer>

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
                  <span>({userName})</span>
                  <span className="ms-2">id-{userId}</span>
                  <span className="ms-2">role-{userRole}</span>
                </Dropdown.Toggle>

                <Dropdown.Menu>
                  <Dropdown.Item href="/profile">
                    <PersonCircle className="me-2" /> Мой профиль
                  </Dropdown.Item>
                  <Dropdown.Item
                    onClick={() => {
                      // Заглушка выхода
                      alert("Выход будет реализован позже");
                    }}
                  >
                    <BoxArrowInRight className="me-2" /> Выйти
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
