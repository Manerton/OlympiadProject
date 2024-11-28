import { useEffect, useState } from "react";
import { Button, Modal } from "react-bootstrap";
import { useParams } from "react-router-dom";
import { MyEvent } from "../../../types/event";
import EventList from "../eventList"


function SubStagePage() {
    const [showModal, setShowModal] = useState(false);
    const [events, setEvents] = useState<MyEvent[]>([]);
    const [event, setEvent] = useState<MyEvent>();

    const { id } = useParams<{ id: string }>();

    useEffect(() => {
        const fetchData = async () => {
            try {
                const [eventResponse, eventsResponse] = await Promise.all([
                    fetch(`http://localhost:8080/events/${id}`, { method: "GET", headers: { "Content-Type": "application/json" } }),
                    fetch(`http://localhost:8080/events/child/${id}`, { method: "GET", headers: { "Content-Type": "application/json" } })
                ]);

                if (!eventResponse.ok || !eventsResponse.ok) {
                    const errorText = await Promise.all([eventResponse.text(), eventsResponse.text()]);
                    throw new Error(`Ошибка API: ${errorText.join(", ")}`);
                }

                const [eventResult, eventsResult] = await Promise.all([
                    eventResponse.json(),
                    eventsResponse.json()
                ]);

                console.log("Этап получен!", eventResult);
                console.log("Подэтапы получены!", eventsResult);

                setEvent(eventResult.data);
                setEvents(eventsResult.data);
            } catch (error) {
                console.error("Ошибка при получении данных:", error);
            }
        };

        fetchData();
    }, [id]);

    return (
        <div>
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
                        </Modal.Body>
                        <Modal.Footer>
                            <Button variant="secondary" onClick={() => setShowModal(false)}>Закрыть</Button>
                        </Modal.Footer>
                    </Modal>
                )}
            </div>
            <div>
                <h6>Дата начала:  {new Date(event?.StartDate || "").toLocaleString()}</h6>
                <h6>Дата конца: {new Date(event?.EndDate || "").toLocaleString()}</h6>


                <h6>Тип: {event?.EventType}</h6>
            </div>
            <div className="d-flex justify-content-between">
                <div>
                    <EventList events={events}/>
                </div>
                <div>
                    
                </div>
            </div>
        </div>

    );
}

export default SubStagePage;
