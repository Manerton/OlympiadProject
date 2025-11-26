import { useEffect, useState } from "react";
import { Table, Button, Form, Spinner, Alert } from "react-bootstrap";
import { useAuth } from "../../Helpers/AuthContext";
import axios from "axios";
import { MyEvent } from "../../types/event";
import { fetchSimpleOlympiads } from "../../../requests/EventsRequests";
import { useParams } from "react-router-dom";

interface Props {
    user_class: number;
    reloadFlag: number;
    onApplied: () => void;
}

const OlympiadsSimpleTable: React.FC<Props> = ({ user_class, onApplied, reloadFlag }) => {
    const { user, accessToken } = useAuth();

    const [olympiads, setOlympiads] = useState<MyEvent[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const [selectedClasses, setSelectedClasses] = useState<Record<string, number>>({});
    const [selectedProfiles, setSelectedProfiles] = useState<Record<string, string>>({});

    useEffect(() => {
        setLoading(true);

       

        fetchSimpleOlympiads()
            .then((res) => setOlympiads(res.data))
            .catch((err) => setError((err as Error).message))
            .finally(() => setLoading(false));
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
            class_number: classNumber,
        };

        try {
            await axios.post(
                "http://localhost:PORT/api/application",
                body,
                { headers: { Authorization: `Bearer ${accessToken}` } }
            );

            alert("Заявка отправлена!");
            onApplied();
        } catch (e) {
            alert("Ошибка при отправке заявки");
        }
    };

    if (loading) return <Spinner />;
    if (error) return <Alert variant="danger">{error}</Alert>;

    return (
        <div className="table-responsive">
            <Table bordered hover className="align-middle text-center">
                <thead>
                    <tr>
                        <th>Предмет</th>
                        <th>Даты проведения</th>
                        <th>Класс участия</th>
                        <th>Профиль</th>
                        <th></th>
                    </tr>
                </thead>

                <tbody>
                    {olympiads.map((olymp) => {

                        const fullDates = olymp.dates.join(", ");

                        return (
                            <tr key={olymp.id}>
                                <td>{olymp.name}</td>

                                <td>{fullDates}</td>

                                <td>
                                    {[9, 10, 11]
                                        .filter((c) => c >= user_class)
                                        .map((c) => (
                                            <Form.Check
                                                key={c}
                                                inline
                                                label={c.toString()}
                                                type="radio"
                                                name={`class-${olymp.id}`}
                                                checked={selectedClasses[olymp.id!] === c}
                                                onChange={() => handleClassChange(olymp.id!, c)}
                                            />
                                        ))}
                                </td>

                                <td>
                                    {olymp.profiles && olymp.profiles.length > 0 ? (
                                        <Form.Select
                                            value={selectedProfiles[olymp.id!] ?? "none"}
                                            onChange={(e) => handleProfileChange(olymp.id!, e.target.value)}
                                        >
                                            <option value="none">Выберите профиль</option>
                                            {olymp.profiles.map((ev) => (
                                                <option key={ev} value={ev}>
                                                    {ev}
                                                </option>
                                            ))}
                                        </Form.Select>
                                    ) : (
                                        "—"
                                    )}
                                </td>

                                <td>
                                    <Button
                                        variant="primary"
                                        className="w-100"
                                        disabled={
                                            !selectedClasses[olymp.id!] || // не выбран класс
                                            (olymp.profiles?.length > 0 &&
                                                (!selectedProfiles[olymp.id!] || selectedProfiles[olymp.id!] === "none")) // есть профили, но не выбран
                                        }
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
        </div>
    );
};

export default OlympiadsSimpleTable;