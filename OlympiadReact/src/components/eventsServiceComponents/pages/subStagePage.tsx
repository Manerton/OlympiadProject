import { useEffect, useState } from "react";
import { Button, Modal, Form } from "react-bootstrap";
import { useParams } from "react-router-dom";
import { MyEvent } from "../../../types/event";
import EventList from "../eventList";
import EventModalForm from "../eventModalWindow";

function SubStagePage() {
    // Триггер показа модального окна
    const [showModal, setShowModal] = useState(false);
    // Список событий 
    const [events, setEvents] = useState<MyEvent[]>([]);
    // Событие
    const [event, setEvent] = useState<MyEvent>();
    // Список жюри
    const [juries, setJuries] = useState<{ id: string; name: string }[]>([]);
    // Выбранные жюри
    const [selectedJuries, setSelectedJuries] = useState<string[]>([]);

    const { id } = useParams<{ id: string }>();

    useEffect(() => {
        const fetchData = async () => {
            try {
                const [eventResponse, eventsResponse, juriesResponse] = await Promise.all([
                    fetch(`http://localhost:8080/events/${id}`, 
                        { 
                            method: "GET", 
                            credentials: "include", // Отправка cookie
                            headers: { "Content-Type": "application/json" }
                        }),
                    fetch(`http://localhost:8080/events/child/${id}`, 
                        { 
                            method: "GET", 
                            credentials: "include", // Отправка cookie
                            headers: { "Content-Type": "application/json" } 
                        }),
                    // TODO! Запрос на UserService для полчения списка жюри
                    fetch(`http://localhost:8080/juries`, 
                        { 
                            method: "GET", 
                            credentials: "include", // Отправка cookie
                            headers: { "Content-Type": "application/json" } 
                        }) // Запрос списка жюри
                ]);

                if (!eventResponse.ok || !eventsResponse.ok || !juriesResponse.ok) {
                    const errorText = await Promise.all([
                        eventResponse.text(),
                        eventsResponse.text(),
                        juriesResponse.text()
                    ]);
                    throw new Error(`Ошибка API: ${errorText.join(", ")}`);
                }

                const [eventResult, eventsResult, juriesResult] = await Promise.all([
                    eventResponse.json(),
                    eventsResponse.json(),
                    juriesResponse.json()
                ]);

                console.log("Этап получен!", eventResult);
                console.log("Подэтапы получены!", eventsResult);
                console.log("Жюри получены!", juriesResult);

                setEvent(eventResult.data);
                setEvents(eventsResult.data);
                setJuries(juriesResult.data); // Устанавливаем список жюри
            } catch (error) {
                console.error("Ошибка при получении данных:", error);
            }
        };

        fetchData();
    }, [id]);

    
    const handleJuryChange = (juryId: string) => {
        setSelectedJuries((prevSelected) =>
            prevSelected.includes(juryId)
                ? prevSelected.filter((id) => id !== juryId)
                : [...prevSelected, juryId]
        );
    };

    // Отправка на сервис juryAssignments для создания множества записей
    const handleSave = async () => {
        try {
            const response = await fetch("http://localhost:8090/jury-assignments-many", {
                method: "POST",
                credentials: "include", // Отправка cookie
                headers: {
                    "Content-Type": "application/json",
                    
                },
                body: JSON.stringify({
                    eventId: id,
                    juryIds: selectedJuries,
                }),
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`Ошибка API: ${errorText}`);
            }

            console.log("Жюри успешно сохранены!");
            alert("Жюри успешно сохранены!");
        } catch (error) {
            console.error("Ошибка при сохранении жюри:", error);
            alert("Произошла ошибка при сохранении жюри.");
        }
    };

    return (
        <div className="container">
            <div className="d-flex justify-content-between">
                <h1>{event?.Name || "Неизвестный этап"}</h1>
                <Button
                    variant="primary"
                    className="mb-3"
                    onClick={() => setShowModal(true)}
                >
                    Создать
                </Button>

                {showModal && (
                    <Modal show onHide={() => setShowModal(false)}>
                        <Modal.Header closeButton>
                            <Modal.Title>Создать подэтап</Modal.Title>
                        </Modal.Header>
                        <Modal.Body>
                            {/* Форма для создания подэтапа */}
                            <EventModalForm event={event} showSubjectField={false} ></EventModalForm>
                        </Modal.Body>
                        <Modal.Footer>
                            <Button variant="secondary" onClick={() => setShowModal(false)}>Закрыть</Button>
                        </Modal.Footer>
                    </Modal>
                )}
            </div>

            <div className="mb-3">
                <h6>Дата начала: {new Date(event?.StartDate || "").toLocaleString()}</h6>
                <h6>Дата конца: {new Date(event?.EndDate || "").toLocaleString()}</h6>
                <h6>Тип: {event?.EventType}</h6>
            </div>

            <div className="row">
                {/* Список подэтапов */}
                <div className="col-8 border p-3">
                    <h5>Список подэтапов</h5>
                    <EventList events={events} />
                </div>

                {/* Список жюри */}
                <div className="col-4 border p-3">
                    <h5>Список жюри</h5>
                    <Button
                        variant="success"
                        className="mb-2"
                        onClick={handleSave}
                    >
                        Сохранить
                    </Button>
                    <Form>
                        {juries.map((jury) => (
                            <Form.Check
                                key={jury.id}
                                type="checkbox"
                                id={`jury-${jury.id}`}
                                label={jury.name}
                                checked={selectedJuries.includes(jury.id)}
                                onChange={() => handleJuryChange(jury.id)}
                            />
                        ))}
                    </Form>
                </div>
            </div>
        </div>
    );
}

export default SubStagePage;
