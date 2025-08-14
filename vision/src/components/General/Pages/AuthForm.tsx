// src/components/AuthForm.tsx
import React, { useState } from "react";
import {
  Container,
  Row,
  Col,
  Form,
  Button,
  InputGroup,
  ListGroup,
} from "react-bootstrap";
import { EyeFill, EyeSlashFill } from "react-bootstrap-icons";

const mockSchools = [
  "Школа №1",
  "Школа №2",
  "Лицей №3",
  "Гимназия №4",
  "Школа-интернат №5",
  "Лицей №6",
];

const AuthForm: React.FC = () => {
  const [isRegister, setIsRegister] = useState(true);
  const [showPassword, setShowPassword] = useState(false);

  // Состояния для поиска школы
  const [schoolQuery, setSchoolQuery] = useState("");
  const [filteredSchools, setFilteredSchools] = useState<string[]>([]);
  const [showList, setShowList] = useState(false);

  const handleSchoolChange = (value: string) => {
    setSchoolQuery(value);
    if (value.trim() === "") {
      setFilteredSchools([]);
      setShowList(false);
      return;
    }
    const results = mockSchools.filter((school) =>
      school.toLowerCase().includes(value.toLowerCase())
    );
    setFilteredSchools(results);
    setShowList(true);
  };

  const selectSchool = (school: string) => {
    setSchoolQuery(school);
    setShowList(false);
  };

  return (
    <Container fluid className="vh-100 d-flex align-items-center">
      <Row className="w-100">
        {/* Левая панель */}
        <Col
          md={6}
          className="d-flex flex-column justify-content-center align-items-center text-center"
          style={{
            background: "linear-gradient(135deg, #3a8dde, #4fd1c5)",
            color: "#fff",
            padding: "2rem",
            borderRadius: "10px 0 0 10px",
          }}
        >
          <h1 className="fw-bold">Добро пожаловать</h1>
          <p>
            Пройдите{" "}
            {isRegister
              ? "регистрацию для участия в ВСОШ"
              : "авторизацию в системе"}
          </p>
        </Col>

        {/* Правая панель */}
        <Col
          md={6}
          className="bg-info-subtle p-4 d-flex flex-column justify-content-center border border-1"
          style={{ borderRadius: "0 10px 10px 0" }}
        >
          {/* Заголовок */}
          <div className="d-flex justify-content-center align-items-center mb-4">
            <h4
              className={`fw-bold m-0 px-2 ${
                isRegister ? "text-primary" : "text-secondary"
              }`}
              style={{ cursor: "pointer" }}
              onClick={() => setIsRegister(true)}
            >
              Регистрация
            </h4>

            <span
              className="mx-2"
              style={{ fontSize: "1.2rem", color: "#999" }}
            >
              /
            </span>

            <h4
              className={`fw-bold m-0 px-2 ${
                !isRegister ? "text-primary" : "text-secondary"
              }`}
              style={{ cursor: "pointer" }}
              onClick={() => setIsRegister(false)}
            >
              Авторизация
            </h4>
          </div>

          {/* Форма */}
          <Form>
            {isRegister && (
              <>
                <Form.Group className="mb-3">
                  <Form.Control type="text" placeholder="Имя" />
                </Form.Group>
                <Form.Group className="mb-3">
                  <Form.Control type="text" placeholder="Фамилия" />
                </Form.Group>
                <Form.Group className="mb-3">
                  <Form.Control type="text" placeholder="Отчество" />
                </Form.Group>

                {/* Поле выбора образовательного учреждения с поиском */}
                <Form.Group className="mb-3 position-relative">
                  <Form.Control
                    type="text"
                    placeholder="Образовательное учреждение"
                    value={schoolQuery}
                    onChange={(e) => handleSchoolChange(e.target.value)}
                    onBlur={() => setTimeout(() => setShowList(false), 150)}
                    onFocus={() => schoolQuery && setShowList(true)}
                  />
                  {showList && filteredSchools.length > 0 && (
                    <ListGroup
                      style={{
                        position: "absolute",
                        top: "100%",
                        width: "100%",
                        zIndex: 1000,
                        maxHeight: "150px",
                        overflowY: "auto",
                      }}
                    >
                      {filteredSchools.map((school, idx) => (
                        <ListGroup.Item
                          key={idx}
                          action
                          onClick={() => selectSchool(school)}
                        >
                          {school}
                        </ListGroup.Item>
                      ))}
                    </ListGroup>
                  )}
                </Form.Group>

                {/* Поле выбора класса с 5 по 11 */}
                <Form.Group className="mb-3">
                  <Form.Select>
                    {[...Array(7)].map((_, i) => {
                      const grade = i + 5;
                      return <option key={grade}>{grade}</option>;
                    })}
                  </Form.Select>
                </Form.Group>
              </>
            )}

            <Form.Group className="mb-3">
              <Form.Control type="email" placeholder="Почта" />
            </Form.Group>

            <Form.Group className="mb-4">
              <InputGroup>
                <Form.Control
                  type={showPassword ? "text" : "password"}
                  placeholder="Пароль"
                />
                <Button
                  variant="outline-secondary"
                  onClick={() => setShowPassword(!showPassword)}
                >
                  {showPassword ? <EyeSlashFill /> : <EyeFill />}
                </Button>
              </InputGroup>
            </Form.Group>

            <Button variant="primary" className="w-100" type="submit">
              {isRegister ? "Создать аккаунт" : "Войти"}
            </Button>
          </Form>
        </Col>
      </Row>
    </Container>
  );
};

export default AuthForm;
