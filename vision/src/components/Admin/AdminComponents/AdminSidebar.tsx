import React, { useState } from "react";
import { Nav, Button, Image } from "react-bootstrap";
import { PersonCircle, BoxArrowInRight, CalendarCheck, PersonBadge } from "react-bootstrap-icons";
import { Link } from "react-router-dom";
import {
  HouseDoor,
  Person,
  FileText,
  Calendar,
  BarChart,
  BoxArrowLeft,
  People,
  FileEarmarkCheck,
  FileEarmarkText,
  CreditCard,
} from "react-bootstrap-icons";

import ThemeToggleButton from "../../Helpers/ThemeToggleButton";

const AdminSidebar: React.FC = () => {
  const [collapsed, setCollapsed] = useState(false);

  const toggleSidebar = () => setCollapsed(!collapsed);

  return (
    <div
      className={`d-flex flex-column justify-content-between bg-dark text-white border-end ${
        collapsed ? "px-2" : "px-3"
      }`}
      style={{
        width: collapsed ? "70px" : "250px",
        height: "100vh",
        position: "sticky",
        top: 0,
        zIndex: 1020, // выше контента, но ниже модалок
        transition: "width 0.3s",
      }}
    >
      {/* Верхняя часть: Аватар и имя */}
      <div>
        <div className="d-flex align-items-center justify-content-between py-3">
          <div className="d-flex align-items-center">
           <PersonCircle size={24} className="me-2" />
            {!collapsed && <strong>Иван Иванович</strong>}
          </div>
          <Button
            variant="link"
            className="text-white p-0 ms-auto"
            onClick={toggleSidebar}
          >
            <span style={{ fontSize: "1.2rem" }}>{collapsed ? "»" : "«"}</span>
          </Button>
        </div>

        <hr className="border-secondary" />

        {/* Навигация */}
        <Nav className="flex-column gap-2">
          <div className="text-uppercase small px-2">
            {!collapsed && "Разделы"}
          </div>

          <Nav.Link as={Link} to="/" className="text-white d-flex align-items-center">
            <HouseDoor className="me-2" />
            {!collapsed && "На главную"}
          </Nav.Link>

          <Nav.Link as={Link} to="/olymp-admin/user/index" className="text-white d-flex align-items-center">
            <People className="me-2" />
            {!collapsed && "Пользователи"}
          </Nav.Link>

          <Nav.Link as={Link} to="/olymp-admin/participant/index" className="text-white d-flex align-items-center">
            <PersonBadge className="me-2" />
            {!collapsed && "Участники"}
          </Nav.Link>

          <Nav.Link as={Link} to="/olymp-admin/school/index" className="text-white d-flex align-items-center">
            <FileEarmarkText className="me-2" />
            {!collapsed && "Школы"}
          </Nav.Link>

          <Nav.Link as={Link} to="/olymp-admin/report/index" className="text-white d-flex align-items-center">
            <CreditCard className="me-2" />
            {!collapsed && "Отчеты"}
          </Nav.Link>

          <Nav.Link as={Link} to="/" className="text-white d-flex align-items-center">
            <CalendarCheck className="me-2" />
            {!collapsed && "Олимпиады"}
          </Nav.Link>

          <Nav.Link as={Link} to="/" className="text-white d-flex align-items-center">
            <FileText className="me-2" />
            {!collapsed && "Заявки"}
          </Nav.Link>

          <Nav.Link as={Link} to="/attendance" className="text-white d-flex align-items-center">
            <Calendar className="me-2" />
            {!collapsed && "Явки"}
          </Nav.Link>

          
        </Nav>
      </div>

      {/* Нижняя часть: Тема и Выход */}
      <div className="pb-3">
        <div className="mb-2 px-2">{<ThemeToggleButton />}</div>
        <Nav.Link as={Link} to="/logout" className="text-danger d-flex align-items-center px-2">
          <BoxArrowLeft className="me-2" />
          {!collapsed && "Выход"}
        </Nav.Link>
      </div>
    </div>
  );
};

export default AdminSidebar;
