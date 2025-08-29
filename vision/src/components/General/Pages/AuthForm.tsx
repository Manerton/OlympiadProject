import React, { useEffect, useRef, useState } from "react";
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
import { axiosSSOAllSchools } from "../../../requests/SSORequests";
import type { School } from "../../types/schools";

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
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [passwordMismatch, setPasswordMismatch] = useState(false);

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
  const [confirmPassword, setConfirmPassword] = useState("");
  const [schoolQuery, setSchoolQuery] = useState("");

  // Состояния для поиска школы
  const [showList, setShowList] = useState(false);

  const [allSchools, setAllSchools] = useState<School[]>([]);
  const [filteredSchools, setFilteredSchools] = useState<School[]>([]);
  const [selectedSchoolId, setSelectedSchoolId] = useState<string>("");


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

    const results = allSchools.filter((school) =>
      school.name.toLowerCase().includes(value.toLowerCase())
    );
    setFilteredSchools(results);
    setShowList(true);
  };

  const selectSchool = (school: School) => {
    setSchoolQuery(school.name);  // показываем название
    setShowList(false);
    // если нужно сохранять id для регистрации
    setSelectedSchoolId(school.id);
  };

  // Проверка совпадения паролей
  const checkPasswordMatch = (pwd: string, confirmPwd: string) => {
    if (isRegister && pwd && confirmPwd) {
      setPasswordMismatch(pwd !== confirmPwd);
    } else {
      setPasswordMismatch(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (isRegister) {
      // Блокируем отправку, если пароли не совпадают
      if (passwordMismatch) {
        return;
      }

      const registerData: RegisterForm = {
        firstname: firstName,
        surname: surName,
        patronymic: patronymic,
        email: email,
        password: password,
        phone_number: phoneNumber,
        gender: gender.toString(),
        school_id: selectedSchoolId,
        birthdate: birthdate,
        classnumber: classNumber.toString(),
        disability: "1",
      };

      try {
        await register(registerData);
      } catch (error) {
        console.error("Registration error:", error);
      }
    } else {
      try {
        await login(email, password);
        navigate("/");
      } catch (error) {
        console.error("Login error:", error);
      }
    }
  };

  useEffect(() => {
    (async () => {
      try {
        const schools = await axiosSSOAllSchools();
        setAllSchools(schools);
      } catch (err) {
        console.error("Ошибка загрузки школ:", err);
      }
    })();
  }, []);

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
              className={`fw-bold m-0 px-2 ${!isRegister ? "text-primary" : "text-secondary"
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
              className={`fw-bold m-0 px-2 ${isRegister ? "text-primary" : "text-secondary"
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
                    <option value={1}>Мужской</option>
                    <option value={2}>Женский</option>
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
                      {filteredSchools.map((school) => (
                        <ListGroup.Item
                          key={school.id}
                          action
                          onClick={() => selectSchool(school)}
                        >
                          {school.name} (регион {school.region})
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

            <Form.Group className="mb-3">
              <InputGroup>
                <Form.Control
                  type={showPassword ? "text" : "password"}
                  placeholder="Пароль"
                  value={password}
                  onChange={(e) => {
                    setPassword(e.target.value);
                    checkPasswordMatch(e.target.value, confirmPassword);
                  }}
                />
                <Button
                  variant="outline-secondary"
                  onClick={() => setShowPassword(!showPassword)}
                >
                  {showPassword ? <EyeSlashFill /> : <EyeFill />}
                </Button>
              </InputGroup>
            </Form.Group>

            {isRegister && (
              <Form.Group className="mb-3">
                <InputGroup>
                  <Form.Control
                    type={showConfirmPassword ? "text" : "password"}
                    placeholder="Повторите пароль"
                    value={confirmPassword}
                    onChange={(e) => {
                      setConfirmPassword(e.target.value);
                      checkPasswordMatch(password, e.target.value);
                    }}
                    isInvalid={passwordMismatch}
                  />
                  <Button
                    variant="outline-secondary"
                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                  >
                    {showConfirmPassword ? <EyeSlashFill /> : <EyeFill />}
                  </Button>
                </InputGroup>
                <Form.Control.Feedback type="invalid">
                  Пароли не совпадают
                </Form.Control.Feedback>
              </Form.Group>
            )}

            <Button
              variant="primary"
              className="w-100"
              type="submit"
              disabled={isRegister && passwordMismatch}
            >
              {isRegister ? "Создать аккаунт" : "Войти"}
            </Button>
          </Form>
        </Col>
      </Row>
    </Container>
  );
};

export default AuthForm;