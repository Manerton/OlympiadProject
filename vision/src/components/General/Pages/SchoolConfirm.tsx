import React, { useState } from "react";
import { Table, Button, Collapse, Card, ListGroup } from "react-bootstrap";

const CITIZENSHIP_TEXT: Record<number, string> = {
    1: "Россия",
    2: "Другое"
};

const DISABILITY_TEXT: Record<number, string> = {
    1: "Нет",
    2: "Есть"
};

// Моковые заявки (5 штук) с категориями и корректными классами
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
    },
    {
        id: "3",
        olympiadName: "Олимпиада по биологии",
        category: "9-11",
        participant: {
            surName: "Кузнецов",
            firstName: "Никита",
            patronymic: "Андреевич",
            birthdate: "2010-11-02",
            email: "kuznetsov@example.com",
            classNumber: "11",
            citizenship: 1,
            disability: 1
        },
        status: "pending"
    },
    {
        id: "4",
        olympiadName: "Олимпиада по истории",
        category: "10",
        participant: {
            surName: "Морозова",
            firstName: "Елизавета",
            patronymic: "Игоревна",
            birthdate: "2007-01-18",
            email: "morozova@example.com",
            classNumber: "10",
            citizenship: 1,
            disability: 1
        },
        status: "pending"
    },
    {
        id: "5",
        olympiadName: "Олимпиада по русскому языку",
        category: "11",
        participant: {
            surName: "Абрамян",
            firstName: "Георгий",
            patronymic: "Самвелович",
            birthdate: "2008-08-08",
            email: "abramyan@example.com",
            classNumber: "11",
            citizenship: 2,
            disability: 1
        },
        status: "pending"
    }
];

const MockApplicationsPage: React.FC = () => {
    const [applications, setApplications] = useState(mockApplications);
    const [openRow, setOpenRow] = useState<string | null>(null);

    const updateStatus = (id: string, newStatus: string) => {
        setApplications(prev =>
            prev.map(app =>
                app.id === id ? { ...app, status: newStatus } : app
            )
        );
    };

    return (
        <div className="container py-4">
            <h2 className="mb-4">Подтверждение заявок учащихся МБОУ г. Астрахани "Лицей № 2"</h2>

            <Table bordered hover responsive>
                <thead>
                <tr>
                    <th>№</th>
                    <th>ФИО ученика</th>
                    <th>Класс ученика</th>
                    <th>Олимпиада</th>
                    <th>Возрастная группа (классов)</th>
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

export default MockApplicationsPage;
