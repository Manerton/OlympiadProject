import React, { useEffect, useState } from "react";
import { Table, Button, Collapse, Card, ListGroup, Alert } from "react-bootstrap";
import { jwtDecode } from "jwt-decode";

// Тип короткоживущего токена
interface JwtPayloadShortLived {
    id: string;
    role: number;
}

// ===== Моки (можно удалить позже) =====

const CITIZENSHIP_TEXT: Record<number, string> = {
    1: "Россия",
    2: "Другое"
};

const DISABILITY_TEXT: Record<number, string> = {
    1: "Нет",
    2: "Есть"
};

const mockApplications = [
    {
        id: "1",
        olympiadName: "Олимпиада по математике",
        category: "9",
        participant: {
            surName: "Иванов",
            firstName: "Петр",
            patronymic: "Сергеевич",
            birthdate: "2008-05-12",
            email: "ivanov@example.com",
            classNumber: "9",
            citizenship: 1,
            disability: 1
        },
        status: "pending"
    },
    {
        id: "2",
        olympiadName: "Олимпиада по физике",
        category: "9-10",
        participant: {
            surName: "Смирнова",
            firstName: "Анна",
            patronymic: "Дмитриевна",
            birthdate: "2009-03-20",
            email: "smirnova@example.com",
            classNumber: "10",
            citizenship: 2,
            disability: 2
        },
        status: "pending"
    }
];

// ===== Компонент =====

const VerifyApplicationsPage: React.FC = () => {
    const [token, setToken] = useState<string | null>(null);
    const [claims, setClaims] = useState<JwtPayloadShortLived | null>(null);
    const [error, setError] = useState<string | null>(null);

    const [applications, setApplications] = useState(mockApplications);
    const [openRow, setOpenRow] = useState<string | null>(null);

    useEffect(() => {
        // 1. Получение токена из URL: /page?access_token=xxxx
        const params = new URLSearchParams(window.location.search);
        const t = params.get("access_token");

        if (!t) {
            setError("Токен доступа отсутствует в URL");
            return;
        }

        setToken(t);

        try {
            // 2. Декодирование JWT
            const decoded = jwtDecode<JwtPayloadShortLived>(t);
            setClaims(decoded);
        } catch (e) {
            setError("Не удалось декодировать токен");
        }
    }, []);

    const updateStatus = (id: string, newStatus: string) => {
        setApplications(prev =>
            prev.map(app =>
                app.id === id ? { ...app, status: newStatus } : app
            )
        );
    };

    // Пример использования токена в headers (который ты позже добавишь)
    const exampleHeaders = token
        ? {
            Authorization: `Bearer ${token}`
        }
        : {};

    // ===== UI =====
    return (
        <div className="container py-4">
            <h2 className="mb-4">Подтверждение заявок учащихся</h2>

            {error && <Alert variant="danger">{error}</Alert>}

            {!token && !error && (
                <Alert variant="warning">Ожидание токена...</Alert>
            )}

            {claims && (
                <Alert variant="info">
                    <strong>Расшифрованные данные токена:</strong><br />
                    ID: {claims.id}<br />
                    Role: {claims.role}
                </Alert>
            )}

            <Table bordered hover responsive>
                <thead>
                <tr>
                    <th>№</th>
                    <th>ФИО ученика</th>
                    <th>Класс</th>
                    <th>Олимпиада</th>
                    <th>Категория</th>
                    <th>Статус</th>
                    <th></th>
                </tr>
                </thead>

                <tbody>
                {applications.map((app, index) => (
                    <React.Fragment key={app.id}>
                        <tr>
                            <td>{index + 1}</td>
                            <td>
                                {app.participant.surName}{" "}
                                {app.participant.firstName}{" "}
                                {app.participant.patronymic}
                            </td>
                            <td>{app.participant.classNumber}</td>
                            <td>{app.olympiadName}</td>
                            <td>{app.category}</td>
                            <td>
                                {app.status === "pending" && "⏳ На рассмотрении"}
                                {app.status === "approved" && "✅ Одобрено"}
                                {app.status === "rejected" && "❌ Отклонено"}
                            </td>
                            <td>
                                <Button
                                    size="sm"
                                    variant="secondary"
                                    onClick={() =>
                                        setOpenRow(openRow === app.id ? null : app.id)
                                    }
                                >
                                    Подробнее
                                </Button>
                            </td>
                        </tr>

                        <tr>
                            <td colSpan={7} className="p-0">
                                <Collapse in={openRow === app.id}>
                                    <div className="p-3 bg-light">
                                        <Card>
                                            <Card.Body>
                                                <h5>Дополнительные сведения</h5>

                                                <ListGroup>
                                                    <ListGroup.Item>
                                                        <strong>ФИО:</strong> {app.participant.surName}{" "}
                                                        {app.participant.firstName}{" "}
                                                        {app.participant.patronymic}
                                                    </ListGroup.Item>
                                                    <ListGroup.Item>
                                                        <strong>Дата рождения:</strong> {app.participant.birthdate}
                                                    </ListGroup.Item>
                                                    <ListGroup.Item>
                                                        <strong>Почта:</strong> {app.participant.email}
                                                    </ListGroup.Item>
                                                    <ListGroup.Item>
                                                        <strong>Класс:</strong> {app.participant.classNumber}
                                                    </ListGroup.Item>
                                                    <ListGroup.Item>
                                                        <strong>Гражданство:</strong>{" "}
                                                        {CITIZENSHIP_TEXT[app.participant.citizenship]}
                                                    </ListGroup.Item>
                                                    <ListGroup.Item>
                                                        <strong>ОВЗ:</strong>{" "}
                                                        {DISABILITY_TEXT[app.participant.disability]}
                                                    </ListGroup.Item>
                                                </ListGroup>

                                                <div className="d-flex gap-2 mt-3">
                                                    <Button
                                                        variant="success"
                                                        onClick={() => updateStatus(app.id, "approved")}
                                                    >
                                                        Одобрить
                                                    </Button>

                                                    <Button
                                                        variant="danger"
                                                        onClick={() => updateStatus(app.id, "rejected")}
                                                    >
                                                        Отклонить
                                                    </Button>
                                                </div>
                                            </Card.Body>
                                        </Card>
                                    </div>
                                </Collapse>
                            </td>
                        </tr>
                    </React.Fragment>
                ))}
                </tbody>
            </Table>
        </div>
    );
};

export default VerifyApplicationsPage;
