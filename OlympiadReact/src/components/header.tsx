import React, { useState, useEffect } from "react";
import { Navbar, Nav, Dropdown, Container } from "react-bootstrap";
import { LinkContainer } from "react-router-bootstrap";
import { PersonCircle, BoxArrowInRight, BoxArrowInLeft } from "react-bootstrap-icons";
import { useRole } from "./RoleContext";
import { useNavigate } from "react-router-dom";




function Header() {
  // Simulating authentication state
  // Check authentication (e.g., from localStorage or API)
  const navigate = useNavigate(); // Initialize the navigation hook
  // Logout function
  const { role, id, name, clearRoleData } = useRole();

  async function handleLogout() {
    try {
      const response = await fetch("http://localhost:8081/logout", {
        method: "POST",
        credentials: "include",
      });
  
      if (!response.ok) throw new Error("Выход не удался");
  
      console.log("Выход выполнен успешно");
      clearRoleData(); // Очищаем контекст
      navigate("/login");
    } catch (error) {
      console.error("Ошибка при выходе:", error);
      alert("Ошибка при выходе. Попробуйте снова.");
    }
  }

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
            {/* тут role == 3 или скорее всего организаторам нужно дать прямую ссылку для отметок*/}
            {role === "3" && (
              <LinkContainer to="/attendance">
               <Nav.Link>Отметить присутствие</Nav.Link>
             </LinkContainer>
            )}
           
          </Nav>
          <Nav className="ms-auto">
            {id ? (
              <Dropdown align="end">
                <Dropdown.Toggle
                  variant="light"
                  id="dropdown-basic"
                  className="d-flex align-items-center border-0 bg-transparent"
                >
                  <PersonCircle size={24} className="me-2" />
                  <span>({name})</span>
                  <span>id-{id}</span>
                  <span>role-{role} </span>
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
