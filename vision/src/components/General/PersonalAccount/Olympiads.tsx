import { useEffect, useState } from "react";
import { Table, Button, Form, Spinner, Alert } from "react-bootstrap";
import { useAuth } from "../../Helpers/AuthContext";
import axios from "axios";
import { MyEvent } from "../../types/event";
import { fetchSimpleOlympiads } from "../../../requests/EventsRequests";

interface Props {
    reloadFlag: number;
    onApplied: () => void;  // ← callback чтобы перезагрузить всю страницу
}


const OlympiadsSimpleTable: React.FC<Props> = ({onApplied, reloadFlag}) => {
    const { user, accessToken } = useAuth();

    const [olympiads, setOlympiads] = useState<MyEvent[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    // per-row selections
    const [selectedClasses, setSelectedClasses] = useState<Record<string, number>>({});
    const [selectedProfiles, setSelectedProfiles] = useState<Record<string, string>>({});

    // загрузка олимпиад
    useEffect(() => {
        setLoading(true);

        fetchSimpleOlympiads(user?.id!) // ← твоя функция, подставь нужные параметры
            .then((res) => {
                setOlympiads(res.data);
            })
            .catch((err) => setError((err as Error).message))
            .finally(() => setLoading(false));

        onApplied()
    }, [reloadFlag]);

    const handleClassChange = (eventId: string, value: number) => {
        setSelectedClasses((prev) => ({ ...prev, [eventId]: value }));
    };

    const handleProfileChange = (eventId: string, value: string) => {
        setSelectedProfiles((prev) => ({ ...prev, [eventId]: value }));
    };

    const handleSubmit = async (eventId: string) => {
        const classNumber = selectedClasses[eventId];

        if (!classNumber) {
            alert("Выберите класс участия");
            return;
        }

        const finalEventId =
            selectedProfiles[eventId] && selectedProfiles[eventId] !== "none"
                ? selectedProfiles[eventId]
                : eventId;

        const body = {
            event_id: finalEventId,
            user_id: user!.id.toString(),
            class_number: classNumber
        };

        try {
            await axios.post(
                "http://localhost:PORT/api/application", // ← подставь свой ApplicationService
                body,
                { headers: { Authorization: `Bearer ${accessToken}` } }
            );

            alert("Заявка отправлена!");
        } catch (e) {
            alert("Ошибка при отправке заявки");
        }
    };

    if (loading) return <Spinner />;
    if (error) return <Alert variant="danger">{error}</Alert>;

    return (
        <Table bordered hover>
            <thead>
                <tr>
                    <th>Предмет</th>
                    <th>Дата проведения</th>
                    <th>Класс</th>
                    <th>Профиль</th>
                    <th></th>
                </tr>
            </thead>

            <tbody>
                {olympiads.map((olymp) => {
                    const start = new Date(olymp.start_date).toLocaleDateString();
                    const end = new Date(olymp.end_date).toLocaleDateString();
                    const fullDate = `${start} — ${end}`;

                    return (
                        <tr key={olymp.id}>
                            <td>{olymp.name}</td>

                            <td>{fullDate}</td>

                            <td style={{ width: 160 }}>
                                <Form.Check
                                    inline
                                    label="9"
                                    type="radio"
                                    name={`class-${olymp.id}`}
                                    checked={selectedClasses[olymp.id!] === 9}
                                    onChange={() => handleClassChange(olymp.id!, 9)}
                                />
                                <Form.Check
                                    inline
                                    label="10"
                                    type="radio"
                                    name={`class-${olymp.id}`}
                                    checked={selectedClasses[olymp.id!] === 10}
                                    onChange={() => handleClassChange(olymp.id!, 10)}
                                />
                                <Form.Check
                                    inline
                                    label="11"
                                    type="radio"
                                    name={`class-${olymp.id}`}
                                    checked={selectedClasses[olymp.id!] === 11}
                                    onChange={() => handleClassChange(olymp.id!, 11)}
                                />
                            </td>

                            <td style={{ width: 200 }}>
                                {olymp.events && olymp.events.length > 0 ? (
                                    <Form.Select
                                        value={selectedProfiles[olymp.id!] ?? "none"}
                                        onChange={(e) => handleProfileChange(olymp.id!, e.target.value)}
                                    >
                                        <option value="none">Нет</option>
                                        {olymp.events.map((ev) => (
                                            <option key={ev.id} value={ev.id}>
                                                {ev.name}
                                            </option>
                                        ))}
                                    </Form.Select>
                                ) : (
                                    "—"
                                )}
                            </td>

                            <td style={{ width: 140 }}>
                                <Button
                                    variant="primary"
                                    className="w-100"
                                    onClick={() => handleSubmit(olymp.id!)}
                                >
                                    Подать заявку
                                </Button>
                            </td>
                        </tr>
                    );
                })}
            </tbody>
        </Table>
    );
};

export default OlympiadsSimpleTable;
