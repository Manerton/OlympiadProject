import React, { useRef, useState } from "react";
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
import { useAuth } from "../../Helpers/AuthContext";
import type { RegisterForm } from "../../types/user";
import { useNavigate } from "react-router-dom";
import { useMask } from "@react-input/mask";

const mockSchools = [
  "Школа №1",
  "Школа №2",
  "Лицей №3",
  "Гимназия №4",
  "Школа-интернат №5",
  "Лицей №6",
];

const AuthForm: React.FC = () => {
  const navigate = useNavigate();
  const [isRegister, setIsRegister] = useState(false);
  const [showPassword, setShowPassword] = useState(false);

  // Поля формы
  const [firstName, setFirstName] = useState("");
  const [surName, setSurname] = useState("");
  const [patronymic, setPatronymic] = useState("");
  const [birthdate, setBirthdate] = useState("");
  const [classNumber, setClassNumber] = useState(0);
  const [gender, setGender] = useState(0);
  const [phoneNumber, setPhoneNumber] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [schoolQuery, setSchoolQuery] = useState("");

  // Состояния для поиска школы
  const [filteredSchools, setFilteredSchools] = useState<string[]>([]);
  const [showList, setShowList] = useState(false);

  // State авторизации
  const { login, register } = useAuth();

  // Define phoneInputRef inside the component
  const phoneInputRef = useMask({
    mask: "+7 (___) ___-__-__",
    replacement: { _: /\d/ },
    showMask: true,
  });

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

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (isRegister) {
      const registerData: RegisterForm = {
        firstname: firstName,
        surname: surName,
        patronymic: patronymic,
        email: email,
        password: password,
        phone_number: phoneNumber,
        gender: gender,
        school: schoolQuery,
        birthdate: birthdate,
        classnumber: classNumber,
      };

      await register(registerData);
    } else {
      await login(email, password);
      navigate("/");
    }
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
                !isRegister ? "text-primary" : "text-secondary"
              }`}
              style={{ cursor: "pointer" }}
              onClick={() => setIsRegister(false)}
            >
              Авторизация
            </h4>
            <span
              className="mx-2"
              style={{ fontSize: "1.2rem", color: "#999" }}
            >
              /
            </span>
            <h4
              className={`fw-bold m-0 px-2 ${
                isRegister ? "text-primary" : "text-secondary"
              }`}
              style={{ cursor: "pointer" }}
              onClick={() => setIsRegister(true)}
            >
              Регистрация
            </h4>
          </div>

          {/* Форма */}
          <Form onSubmit={handleSubmit}>
            {isRegister && (
              <>
                <Form.Group className="mb-3">
                  <Form.Control
                    type="text"
                    placeholder="Имя"
                    value={firstName}
                    onChange={(e) => setFirstName(e.target.value)}
                  />
                </Form.Group>
                <Form.Group className="mb-3">
                  <Form.Control
                    type="text"
                    placeholder="Фамилия"
                    value={surName}
                    onChange={(e) => setSurname(e.target.value)}
                  />
                </Form.Group>
                <Form.Group className="mb-3">
                  <Form.Control
                    type="text"
                    placeholder="Отчество"
                    value={patronymic}
                    onChange={(e) => setPatronymic(e.target.value)}
                  />
                </Form.Group>

                {/* Дата рождения */}
                <Form.Group className="mb-3">
                  <Form.Label>Дата рождения</Form.Label>
                  <Form.Control
                    type="date"
                    value={birthdate}
                    onChange={(e) => setBirthdate(e.target.value)}
                    isInvalid={birthdate !== "" && new Date(birthdate) > new Date()}
                  />
                  <Form.Control.Feedback type="invalid">
                    Дата рождения не может быть в будущем
                  </Form.Control.Feedback>
                </Form.Group>

                {/* Пол */}
                <Form.Group className="mb-3">
                  <Form.Label>Пол</Form.Label>
                  <Form.Select
                    value={gender}
                    onChange={(e) => setGender(Number(e.target.value))}
                  >
                    <option value={0}>Мужской</option>
                    <option value={1}>Женский</option>
                  </Form.Select>
                </Form.Group>

                {/* Телефон */}
                <Form.Group className="mb-3">
                  <Form.Control
                    type="tel"
                    placeholder="+7 (___) ___-__-__"
                    ref={phoneInputRef}
                    value={phoneNumber}
                    onChange={(e) => setPhoneNumber(e.target.value)}
                  />
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

                {/* Поле выбора класса */}
                <Form.Group className="mb-3">
                  <Form.Label>Класс</Form.Label>
                  <Form.Select
                    value={classNumber}
                    onChange={(e) => setClassNumber(Number(e.target.value))}
                  >
                    {[...Array(7)].map((_, i) => {
                      const grade = i + 5;
                      return <option key={grade}>{grade}</option>;
                    })}
                  </Form.Select>
                </Form.Group>
              </>
            )}

            <Form.Group className="mb-3">
              <Form.Control
                type="email"
                placeholder="Почта"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
            </Form.Group>

            <Form.Group className="mb-4">
              <InputGroup>
                <Form.Control
                  type={showPassword ? "text" : "password"}
                  placeholder="Пароль"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
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