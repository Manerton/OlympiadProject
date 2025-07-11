import { useEffect, useState } from "react";
import { Button, Modal, Form } from "react-bootstrap";
import { MyEvent } from "../../../types/event";
import EventList from "../eventList";
import EventModalForm from "../eventModalWindow";
import API_CONFIG from "../../../config/apiConfig";
import { useRole } from "../../RoleContext";
import EventInfo from "../eventInfo";
import UserRoles from "../../../types/user";

interface BaseEventPageProps {
    eventId?: string
}

function SubStagMainPage({ eventId }: BaseEventPageProps) {
    // Триггер показа модального окна
    const [showModal, setShowModal] = useState(false);
    // Список событий 
    const [events, setEvents] = useState<MyEvent[]>([]);
    // Событие
    const [event, setEvent] = useState<MyEvent | null>(null); 
    // Список жюри
    const [juries, setJuries] = useState<{ id: string; name: string , role: string}[]>([]);
    // Выбранные жюри
    const [selectedJuries, setSelectedJuries] = useState<string[]>([]);

    const [isLoading, setIsLoading] = useState(true); // Добавляем состояние загрузки
    const [error, setError] = useState<string | null>(null); // Добавляем состояние ошибки


    const { role, id } = useRole();

    const handleDeleteEvent = (id: number) => {
        setEvents((events) => events.filter((event) => event.ID !== id));
    };

    const fetchData = async () => {
        setIsLoading(true);
        setError(null);
        try {
            const [eventResponse, eventsResponse] = await Promise.all([
                fetch(`${API_CONFIG.EVENTS}/${eventId}`,
                    {
                        method: "GET",
                        credentials: "include", // Отправка cookie
                        headers: { "Content-Type": "application/json" }
                    }),
                fetch(`${API_CONFIG.EVENTS}/child/${eventId}`,
                    {
                        method: "GET",
                        credentials: "include", // Отправка cookie
                        headers: { "Content-Type": "application/json" }
                    }),
                // TODO! Запрос на UserService для полчения списка жюри
                // fetch(`http://localhost:8081/juries`, 
                //     { 
                //         method: "GET", 
                //         credentials: "include", // Отправка cookie
                //         headers: { "Content-Type": "application/json" } 
                //     }) // Запрос списка жюри
            ]);

            if (!eventResponse.ok || !eventsResponse.ok ) {
                const errorText = await Promise.all([
                    eventResponse.text(),
                    eventsResponse.text(),
                ]);
                throw new Error(`Ошибка API: ${errorText.join(", ")}`);
            }

            const [eventResult, eventsResult] = await Promise.all([
                eventResponse.json(),
                eventsResponse.json(),
            ]);

            console.log("Этап получен!", eventResult);
            console.log("Подэтапы получены!", eventsResult);

            setEvent(eventResult.data);
            setEvents(eventsResult.data);
            // setJuries(juriesResult.data);
            console.log("Жюри dscf!", juries);

        } catch (error) {
            console.error("Ошибка при получении данных:", error);
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
    }, [eventId]);

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
            const response = await fetch(`${API_CONFIG.JUREASSIGNMENTS}/many`, {
                method: "POST",
                credentials: "include", // Отправка cookie
                headers: {"Content-Type": "application/json",},
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

    const OnUpdateListEvent = () => {
        setShowModal(false)
        fetchData()
      }

    if (isLoading) {
        return <div className="container text-center mt-5">
            <div className="spinner-border" role="status">
                <span className="visually-hidden">Загрузка...</span>
            </div>
        </div>;
    }

    if (error) {
        return <div className="container alert alert-danger mt-3">{error}</div>;
    }

    if (!event) {
        return <div className="container alert alert-warning mt-3">Событие не найдено</div>;
    }

      
    return (
        <div className="container">
            <div className="row">
                <h1>{event?.name || "Неизвестный этап"}</h1>
                <div className="col-3">
                    <h3>Информация о {event?.name}</h3>
                    {event && (
                        <EventInfo event={event}></EventInfo>
                    )}
                </div>

                {/* Список подэтапов */}
                <div className="col-5">
                    <div className="d-flex align-items-center justify-content-between">
                        <h5>Список стадий</h5>
                        {role === UserRoles.Organaizer && (
                            <Button
                                variant="primary"
                                className="mb-2"
                                onClick={() => setShowModal(true)}
                            >
                                Создать
                            </Button>
                        )}

                        {showModal && (
                            <Modal show onHide={() => setShowModal(false)}>
                                <Modal.Header closeButton>
                                    <Modal.Title>Создать стадию</Modal.Title>
                                </Modal.Header>
                                <Modal.Body>
                                    {/* Форма для создания подэтапа */}
                                    <EventModalForm onSuccess={OnUpdateListEvent} event={event} showSubjectField={false} ></EventModalForm>
                                </Modal.Body>
                                <Modal.Footer>
                                    <Button variant="secondary" onClick={() => setShowModal(false)}>Закрыть</Button>
                                </Modal.Footer>
                            </Modal>
                        )}
                    </div>
                    <div className="card shadow-sm mb-2">
                        <div className="card-body">
                            <EventList events={events}  onDelete={handleDeleteEvent}/>
                        </div>
                    </div>
                </div>

                {/* Список жюри */}
                <div className="col-4">
                    <div className="d-flex align-items-center justify-content-between">
                        <h5>Список жюри</h5>

                        {role === UserRoles.Organaizer && (
                        <Button
                            variant="success"
                            className="mb-2"
                            onClick={handleSave}
                        >
                            Сохранить
                        </Button>
                        )}
                    </div>
                    
                    <div className="card shadow-sm mb-2">
                        <div className=" card-body">
                            <Form>
                                {juries.length > 0 ? ( juries.map((jury) => (
                                    <Form.Check className="d-flex justify-content-between "
                                        key={jury.id}
                                        type="checkbox"
                                        id={`jury-${jury.id}`}
                                        label={jury.name}
                                        checked={selectedJuries.includes(jury.id)}
                                        onChange={() => handleJuryChange(jury.id)}
                                    />
                                ))) : (
                                    <p>Список пуст...</p>
                                )}
                            </Form>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default SubStagMainPage;
