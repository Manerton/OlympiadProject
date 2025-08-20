import React from "react";
import { Navbar, Nav, Dropdown, Container } from "react-bootstrap";
import { LinkContainer } from "react-router-bootstrap";
import { Link } from "react-router-dom";
import { PersonCircle, BoxArrowInRight } from "react-bootstrap-icons";
import ThemeToggleButton from '../../Helpers/ThemeToggleButton';
import RegionalStages from '../Pages/Events/RegionalStages';


function Header() {
  // Заглушки: в будущем подключите контекст пользователя
  const isAuthenticated = false; // временная проверка
  const userName = "Имя";
  const userId = "123";
  const userRole = "admin"; // временно установлено 'admin' для демонстрации

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

            <LinkContainer to={`/applications/user/${userId}`}>
              <Nav.Link>Статус заявки</Nav.Link>
            </LinkContainer>

            {/* Админские ссылки, показываются только для админов */}
            {userRole === 'admin' && (
              <>
                <LinkContainer to="/AdminPanel">
                  <Nav.Link>Панель Администрирования</Nav.Link>
                </LinkContainer>
                {/* <Dropdown>
                  <Dropdown.Toggle as={Nav.Link} variant="light">
                    Администрирование
                  </Dropdown.Toggle>
                  <Dropdown.Menu>
                    <LinkContainer to="/olymp-admin/user/index">
                      <Dropdown.Item>Пользователи</Dropdown.Item>
                    </LinkContainer>
                    <LinkContainer to="/olymp-admin/participant/index">
                      <Dropdown.Item>Участники</Dropdown.Item>
                    </LinkContainer>
                    <LinkContainer to="/olymp-admin/school/index">
                      <Dropdown.Item>Школы</Dropdown.Item>
                    </LinkContainer>
                    <LinkContainer to="/olymp-admin/report/index">
                      <Dropdown.Item>Отчёты</Dropdown.Item>
                    </LinkContainer>
                  </Dropdown.Menu>
                </Dropdown> */}
              </>
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
                  <span>({userName})</span>
                  <span className="ms-2">id-{userId}</span>
                  <span className="ms-2">role-{userRole}</span>
                </Dropdown.Toggle>

                <Dropdown.Menu>
                  <Dropdown.Item as={Link} to="/profile">
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