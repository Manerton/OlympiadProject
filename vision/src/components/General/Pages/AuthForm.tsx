import React, { useEffect, useState } from "react";
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
import { useMask } from "@react-input/mask";
import { useAuth } from "../../Helpers/AuthContext";
import type { School } from "../../types/schools";
import type { RegisterForm } from "../../types/user";
import axios from "axios";
import { axiosSendSMSCode } from "../../../requests/NotificationRequests";
import { axiosSSOVerifySMSCode } from "../../../requests/SSORequests";

axios.defaults.baseURL = "http://localhost:6611";

const AuthForm: React.FC = () => {
    const { login, register } = useAuth();

    const [isRegister, setIsRegister] = useState(false);
    const [step, setStep] = useState(1);

    // Шаг 1 — ФИО + дата + пол
    const [firstName, setFirstName] = useState("");
    const [surName, setSurName] = useState("");
    const [patronymic, setPatronymic] = useState("");
    const [birthdate, setBirthdate] = useState("");
    const [gender, setGender] = useState(1);


    // Финальный шаг — SMS подтверждение
    const [phoneNumber, setPhoneNumber] = useState("");
    const [smsSent, setSmsSent] = useState(false);
    const [smsCode, setSmsCode] = useState("");
    const [isPhoneVerified, setIsPhoneVerified] = useState(false);

    // Расширенная логика тайм-аутов
    const [attempts, setAttempts] = useState(0);
    const retryTimeouts = [60, 300, 300]; // 1 мин → 5 мин → 5 мин
    const [cooldown, setCooldown] = useState(0);
    const [errorSMS, setErrorSMS] = useState("");

    const phoneInputRef = useMask({
        mask: "+7 (___) ___-__-__",
        replacement: { _: /\d/ },
        showMask: true,
    });

    // Шаг 3 — школа + класс
    // Districts
    const [districts, setDistricts] = useState([]);
    const [selectedDistrictId, setSelectedDistrictId] = useState("");

    // Шаг 3 - гражданство + овз
    const [citizenship, setCitizenship] = useState(0);
    const [disability, setDisability] = useState(0);



    // Schools
    const [schoolsInDistrict, setSchoolsInDistrict] = useState([]);
    const [selectedSchoolId, setSelectedSchoolId] = useState("");
    const [schoolQuery, setSchoolQuery] = useState("");
    const [classNumber, setClassNumber] = useState(5);

    // Шаг 4 — пароль
    const [email, setEmail] = useState("");
    const [emailError, setEmailError] = useState("");
    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");
    const [showPassword, setShowPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);
    const passwordMismatch = password !== confirmPassword;

    const CITIZENSHIP_TEXT: Record<number, string> = {
        0: "Не выбрано",
        1: "Россия",
        2: "Другое",
    };

    const DISABILITY_TEXT: Record<number, string> = {
        0: "Не выбрано",
        1: "Нет",
        2: "Есть",
    };

// Использование:
    const citizenshipText = CITIZENSHIP_TEXT[citizenship] ?? "Не выбрано";
    const disabilityText = DISABILITY_TEXT[disability] ?? "Не выбрано";


    ////////////////////////////////////////////////
    // SMS Logics
    ////////////////////////////////////////////////

    const sendSMS = async () => {
        try {
            setErrorSMS("");

            const data = await axiosSendSMSCode(phoneNumber);
            // data можно использовать, если backend вернёт что-то вроде { success: true }

            setSmsSent(true);

            // выставляем задержку на повтор
            const timeout = retryTimeouts[Math.min(attempts, retryTimeouts.length - 1)];
            setCooldown(timeout);

        } catch (e) {
            console.error("Error sending SMS:", e);
            setErrorSMS("Ошибка отправки SMS. Попробуйте позже.");
        }
    };


    const verifySMS = async () => {
        try {
            setErrorSMS("");

            const result = await axiosSSOVerifySMSCode(phoneNumber, smsCode);

            if (!result.data.success) {
                throw new Error(result.message || "Неверный код");
            }

            setIsPhoneVerified(true);
            await handleFinalSubmit();

        } catch (e) {
            console.error("Неверный код SMS:", e);

            setErrorSMS("Неверный код. Попробуйте снова.");

            setAttempts(a => {
                const next = a + 1;
                const timeout = retryTimeouts[Math.min(next, retryTimeouts.length - 1)];
                setCooldown(timeout);
                return next;
            });
        }
    };





    useEffect(() => {
        if (cooldown <= 0) return;
        const t = setInterval(() => setCooldown((c) => c - 1), 1000);
        return () => clearInterval(t);
    }, [cooldown]);


    ////////////////////////////////////////////////
    // School Search
    ////////////////////////////////////////////////

    useEffect(() => {
        const fetchDistricts = async () => {
            try {
                const res = await axios.get(`/districts/30`);
                setDistricts(res.data.data);
            } catch (err) {
                console.error("Ошибка загрузки округов:", err);
            }
        };

        fetchDistricts();
    }, []);

    useEffect(() => {
        if (!selectedDistrictId) return;

        const fetchSchools = async () => {
            try {
                const res = await axios.get(`/schools/district/${selectedDistrictId}`);
                setSchoolsInDistrict(res.data.data);
            } catch (err) {
                console.error("Ошибка загрузки школ:", err);
            }
        };

        fetchSchools();
    }, [selectedDistrictId]);


    const validateEmail = (value: string) => {
        const regex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return regex.test(value);
    };

    ////////////////////////////////////////////////
    // Submit
    ////////////////////////////////////////////////

    const handleFinalSubmit = async () => {
        const registerData: RegisterForm = {
            firstname: firstName,
            surname: surName,
            patronymic,
            email,
            password,
            phone_number: phoneNumber,
            gender: gender.toString(),
            school_id: selectedSchoolId,
            birthdate,
            classnumber: classNumber.toString(),
            disability: String(disability),
            citizenship: String(citizenship)
        };

        try {
            await register(registerData);
        } catch (err) {
            console.error("Registration error:", err);
        }
    };

    ////////////////////////////////////////////////
    // Step Components
    ////////////////////////////////////////////////

    const Step1 = () => (
        <>
            <h4 className="fw-bold mb-3">Шаг 1: Личная информация</h4>

            <Form.Group className="mb-3">
                <Form.Control
                    placeholder="Имя"
                    value={firstName}
                    onChange={(e) => setFirstName(e.target.value)}
                />
            </Form.Group>

            <Form.Group className="mb-3">
                <Form.Control
                    placeholder="Фамилия"
                    value={surName}
                    onChange={(e) => setSurName(e.target.value)}
                />
            </Form.Group>

            <Form.Group className="mb-3">
                <Form.Control
                    placeholder="Отчество"
                    value={patronymic}
                    onChange={(e) => setPatronymic(e.target.value)}
                />
            </Form.Group>

            <Form.Group className="mb-3">
                <Form.Label>Дата рождения</Form.Label>
                <Form.Control
                    type="date"
                    value={birthdate}
                    onChange={(e) => setBirthdate(e.target.value)}
                />
            </Form.Group>

            <Form.Group className="mb-3">
                <Form.Label>Пол</Form.Label>
                <Form.Select value={gender} onChange={(e) => setGender(Number(e.target.value))}>
                    <option value={1}>Мужской</option>
                    <option value={2}>Женский</option>
                </Form.Select>
            </Form.Group>

            <Button className="w-100" onClick={() => setStep(2)}>
                Далее
            </Button>
        </>
    );


    const Step2 = () => (
        <>
            <h4 className="fw-bold mb-3">Шаг 2: Образовательное учреждение</h4>

            {/* Выбор муниципального округа */}
            <Form.Group className="mb-3">
                <Form.Label>Муниципальное образование</Form.Label>
                <Form.Select
                    value={selectedDistrictId}
                    onChange={(e) => {
                        setSelectedDistrictId(e.target.value);
                        setSelectedSchoolId(""); // сбрасываем школу
                    }}
                >
                    <option value="">Выберите округ...</option>
                    {districts.map((d: any) => (
                        <option key={d.id} value={d.id}>
                            {d.name}
                        </option>
                    ))}
                </Form.Select>
            </Form.Group>

            {/* Выбор школы */}
            <Form.Group className="mb-3">
                <Form.Label>Школа</Form.Label>
                <Form.Select
                    value={selectedSchoolId}
                    disabled={!selectedDistrictId}
                    onChange={(e) => {
                        setSelectedSchoolId(e.target.value);
                        setSchoolQuery(e.target.name);
                    }}
                >
                    <option value="">Выберите школу...</option>
                    {schoolsInDistrict.map((s: any) => (
                        <option key={s.id} value={s.id}>
                            {s.name}
                        </option>
                    ))}
                </Form.Select>
            </Form.Group>

            {/* Выбор класса */}
            <Form.Group className="mb-3">
                <Form.Label>Класс</Form.Label>
                <Form.Select value={classNumber} onChange={(e) => setClassNumber(Number(e.target.value))}>
                    {[...Array(7)].map((_, idx) => {
                        const n = idx + 5;
                        return (
                            <option key={n} value={n}>
                                {n}
                            </option>
                        );
                    })}
                </Form.Select>
            </Form.Group>

            <Button
                className="w-100"
                disabled={!selectedDistrictId || !selectedSchoolId}
                onClick={() => setStep(3)}
            >
                Далее
            </Button>

            <Button variant="secondary" className="w-100 mt-3" onClick={() => setStep(1)}>
                Назад
            </Button>
        </>
    );

    const Step3 = () => (
        <>
            <h4 className="fw-bold mb-3">Шаг 3: Дополнительная информация</h4>

            {/* Гражданство */}
            <Form.Group className="mb-4">
                <Form.Label className="fw-semibold">
                    Гражданство
                </Form.Label>

                <div className="d-flex flex-column gap-2 mt-2">

                    <Form.Check
                        type="radio"
                        id="citizenship-russia"
                        label="Россия"
                        value="russia"
                        checked={citizenship === 1}
                        onChange={() => setCitizenship(1)}
                    />

                    <Form.Check
                        type="radio"
                        id="citizenship-other"
                        label="Другое"
                        value="other"
                        checked={citizenship === 2}
                        onChange={() => setCitizenship(2)}
                    />
                </div>
            </Form.Group>

            {/* ОВЗ */}
            <Form.Group className="mb-4">
                <Form.Label className="fw-semibold">
                    Наличие ограничений возможностей здоровья (ОВЗ)
                </Form.Label>

                <div className="d-flex flex-column gap-2 mt-2">

                    <Form.Check
                        type="radio"
                        id="disability-no"
                        label="Нет"
                        value="no"
                        checked={disability === 1}
                        onChange={() => setDisability(1)}
                    />

                    <Form.Check
                        type="radio"
                        id="disability-yes"
                        label="Есть"
                        value="yes"
                        checked={disability === 2}
                        onChange={() => setDisability(2)}
                    />
                </div>

                {/* Пояснение для школьников */}
                <Form.Text className="text-muted">
                    ОВЗ — это особенности здоровья, при которых может понадобиться помощь или особые условия для прохождения олимпиады.
                    Если сомневаетесь, лучше выбрать «Нет».
                </Form.Text>
            </Form.Group>

            <Button className="w-100" onClick={() => setStep(4)}>
                Далее
            </Button>
        </>
    );


    const Step4 = () => (
        <>
            <h4 className="fw-bold mb-3">Шаг 4: Данные аккаунта</h4>

            <Form.Group className="mb-3">
                <Form.Control
                    placeholder="Почта"
                    value={email}
                    onChange={(e) => {
                        const value = e.target.value;
                        setEmail(value);

                        if (!validateEmail(value)) {
                            setEmailError("Некорректный формат почты");
                        } else {
                            setEmailError("");
                        }
                    }}
                    isInvalid={!!emailError}
                />

                <Form.Control.Feedback type="invalid">
                    {emailError}
                </Form.Control.Feedback>
            </Form.Group>


            <Form.Group className="mb-3">
                <InputGroup>
                    <Form.Control
                        type={showPassword ? "text" : "password"}
                        placeholder="Пароль"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                    />
                    <Button onClick={() => setShowPassword(!showPassword)}>
                        {showPassword ? <EyeSlashFill /> : <EyeFill />}
                    </Button>
                </InputGroup>
            </Form.Group>

            <Form.Group className="mb-3">
                <InputGroup>
                    <Form.Control
                        type={showConfirmPassword ? "text" : "password"}
                        placeholder="Повторите пароль"
                        value={confirmPassword}
                        isInvalid={passwordMismatch}
                        onChange={(e) => setConfirmPassword(e.target.value)}
                    />
                    <Button onClick={() => setShowConfirmPassword(!showConfirmPassword)}>
                        {showConfirmPassword ? <EyeSlashFill /> : <EyeFill />}
                    </Button>
                </InputGroup>

                {passwordMismatch && (
                    <div className="text-danger mt-1">Пароли не совпадают</div>
                )}

            </Form.Group>

            <Button
                className="w-100"
                disabled={!!emailError || !email || passwordMismatch}
                onClick={() => setStep(5)}
            >
                Далее
            </Button>

            <Button variant="secondary" className="w-100 mt-3" onClick={() => setStep(3)}>
                Назад
            </Button>
        </>
    );

    // А в Step5 используй так:
    const Step5 = () => {
        const selectedSchool = schoolsInDistrict.find((s: any) => s.id === selectedSchoolId);
        const schoolName = selectedSchool?.name || "Не выбрана";

        return (
            <>
                <h4 className="fw-bold mb-3">Шаг 5: Подтверждение</h4>
                <p>Проверьте данные перед завершением регистрации:</p>

                <ListGroup className="mb-4">
                    <ListGroup.Item>ФИО: {surName} {firstName} {patronymic}</ListGroup.Item>
                    <ListGroup.Item>Дата рождения: {birthdate}</ListGroup.Item>
                    <ListGroup.Item>Почта: {email}</ListGroup.Item>
                    <ListGroup.Item>Школа: {schoolName}</ListGroup.Item>
                    <ListGroup.Item>Класс: {classNumber}</ListGroup.Item>
                    <ListGroup.Item>Гражданство: <strong>{CITIZENSHIP_TEXT[citizenship]}</strong></ListGroup.Item>
                    <ListGroup.Item>ОВЗ: <strong>{DISABILITY_TEXT[disability]}</strong></ListGroup.Item>
                </ListGroup>

                <Button className="w-100" onClick={() => setStep(6)}>
                    Всё верно, продолжить
                </Button>
                <Button variant="secondary" className="w-100 mt-3" onClick={() => setStep(4)}>
                    Назад
                </Button>
            </>
        );
    };

    const Step6 = () => (
        <>
            <h4 className="fw-bold mb-3">Почти всё готово!</h4>
            <p>Осталось лишь подтвердить ваш номер телефона:<br /><b>{phoneNumber}</b></p>

            <Form.Group> <Form.Control placeholder="+7 (___) ___-__-__" ref={phoneInputRef} value={phoneNumber} onChange={(e) => setPhoneNumber(e.target.value)} /> </Form.Group>

            {!smsSent && (
                <Button className="w-100 mt-3" onClick={sendSMS}>
                    Отправить SMS-код
                </Button>
            )}

            {smsSent && (
                <>
                    <Form.Control
                        className="mt-3"
                        placeholder="Введите код"
                        value={smsCode}
                        onChange={(e) => setSmsCode(e.target.value)}
                    />

                    {errorSMS && (
                        <div className="text-danger mt-2">{errorSMS}</div>
                    )}

                    <Button className="w-100 mt-3" onClick={verifySMS}>
                        Подтвердить
                    </Button>

                    <div className="mt-3 text-center">
                        {cooldown > 0 ? (
                            <span>Повторная отправка через {cooldown} сек.</span>
                        ) : (
                            <Button
                                variant="link"
                                className="p-0"
                                onClick={sendSMS}
                            >
                                Отправить код повторно
                            </Button>
                        )}
                    </div>
                </>
            )}

            <Button variant="secondary" className="w-100 mt-4" onClick={() => setStep(5)}>
                Назад
            </Button>
        </>
    );

    ////////////////////////////////////////////////
    // Render
    ////////////////////////////////////////////////

    const renderStep = () => {
        switch (step) {
            case 1: return <Step1 />;
            case 2: return <Step2 />;
            case 3: return <Step3 />;
            case 4: return <Step4 />;
            case 5: return <Step5 />;
            case 6: return <Step6 />;
            default: return <Step1 />;
        }
    };

    return (
        <Container fluid className="vh-100 d-flex align-items-center">
            <Row className="w-100">
                <Col md={6} className="d-flex flex-column justify-content-center align-items-center text-center"
                    style={{
                        background: "linear-gradient(135deg, #3a8dde, #4fd1c5)",
                        color: "#fff",
                        padding: "2rem",
                        borderRadius: "10px 0 0 10px",
                    }}
                >
                    <h1 className="fw-bold">Добро пожаловать</h1>
                    <p>{isRegister ? "Регистрация в системе" : "Авторизация"}</p>
                </Col>

                <Col md={6} className="p-4 border border-1 d-flex flex-column"
                    style={{ borderRadius: "0 10px 10px 0", height: "572px" }}
                >
                    {!isRegister ? (
                        <>
                            <h4 className="fw-bold mb-3">Авторизация</h4>

                            <Form.Group className="mb-3">
                                <Form.Control
                                    placeholder="Почта"
                                    value={email}
                                    onChange={(e) => setEmail(e.target.value)}
                                />
                            </Form.Group>

                            <Form.Group className="mb-3">
                                <Form.Control
                                    type="password"
                                    placeholder="Пароль"
                                    value={password}
                                    onChange={(e) => setPassword(e.target.value)}
                                />
                            </Form.Group>

                            <Button className="w-100" onClick={() => login(email, password)}>
                                Войти
                            </Button>

                            <div className="text-center mt-3">
                                <Button variant="link" onClick={() => setIsRegister(true)}>
                                    Создать аккаунт
                                </Button>
                            </div>
                        </>
                    ) : (
                        <>
                            <div className="d-flex justify-content-between mb-4">
                                <h4 className="fw-bold m-0">Регистрация</h4>
                                <span>Шаг {step} / 5</span>
                            </div>

                            {renderStep()}

                            <div className="text-center mt-3">
                                <Button variant="link" onClick={() => setIsRegister(false)}>
                                    Уже есть аккаунт?
                                </Button>
                            </div>
                        </>
                    )}
                </Col>
            </Row>
        </Container>
    );
};

export default AuthForm;
