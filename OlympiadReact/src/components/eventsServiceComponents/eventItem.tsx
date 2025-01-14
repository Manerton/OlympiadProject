import { MyEvent, REGIONAL_STAGE, STAGE, OLYMPIAD, APPEAL, VIEW_WORKS } from "../../types/event";
import { Button, Card } from "react-bootstrap";
import { useNavigate, Link } from "react-router-dom";
import { useRole } from "../RoleContext";
// import { FaTrash, FaChevronDown } from "react-icons/fa";
import API_CONFIG from "../../config/apiConfig";
import UserRoles from "../../types/user";

interface EventItemProps {
  event: MyEvent;
  onDelete: (id: number) => void;
  isSubmitApplication: boolean;
}

function EventItem({ event, onDelete, isSubmitApplication }: EventItemProps) {
  const navigate = useNavigate();
  const { role, id } = useRole();

  const handleClick = () => {
    if (event.EventType === REGIONAL_STAGE) {
      navigate(`/olympiads/${event.ID}`);
    } else if (event.EventType === OLYMPIAD) {
      navigate(`/olympiad-stages/${event.ID}`);
    } else if (event.EventType === STAGE) {
      navigate(`/sub-stage/${event.ID}`)
    }
  }

  const deleteEvent = async () => {

    const endPointDeleteEvent = `${API_CONFIG.EVENTS}/${event.ID}`;
    try {

      const response = await fetch(endPointDeleteEvent, {
        method: 'DELETE',
        credentials: "include", // Отправка cookie
        headers: { "Content-Type": "application/json" }
      });
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const result = await response.json();
      if (result.status === "Error") {
        alert("Это событие связано с другими событиями, поэтому его нельзя удалить!");
      } else {
        onDelete(event.ID || 0)
      }
      console.log("Response from API:", result);
    } catch (error) {
      console.error("Ошибка при загрузке региональных этапов:", error);
    }
  }

  const onSubmitApplication = async () => {
    const endPointSubmitApplication = `${API_CONFIG.APPLICATION}`
    try {
      const response = await fetch(endPointSubmitApplication, {
        method: "POST",
        credentials: "include", // Отправка cookie
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ "user_id": id, "event_id": event.ID }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const result = await response.json();
      console.log("Response from API:", result);

    } catch (error) {
      console.error("Ошибка при подачи заявки")
    }
  }

  return (
    <Card key={event.ID} className="mb-3">
      <Card.Body className="d-flex justify-content-between align-items-start">
        {/* Левая часть: текст названия и даты */}
        <div onClick={handleClick} style={{ cursor: "pointer", textAlign: "left" }}>
          <Card.Title className="fw-bold mb-2" style={{ fontSize: "1.2rem" }}>
            {event.Name}
          </Card.Title>
          {(event.EventType === APPEAL || event.EventType === VIEW_WORKS) && (
            <Card.Title className="fw-bold mb-2" style={{ fontSize: "1.2rem" }}>
              {event.EventType === APPEAL ? ("Апелляция") : ("Просмотр работ")}
            </Card.Title>
          )}

          <Card.Text className="text-muted mb-1">
            <strong>Дата начала:</strong> {new Date(event.StartDate).toLocaleString()}
          </Card.Text>
          <Card.Text className="text-muted mb-1">
            <strong>Дата конца:</strong> {new Date(event.EndDate).toLocaleString()}
          </Card.Text>
          {event.Subject && (
            <Card.Text className="text-muted mb-1">
              <strong>Предмет:</strong> {event.Subject}
            </Card.Text>
          )}
          {event.AdditionalInfo && (
            <Card.Text className="mt-3">
              <strong>Дополнительно:</strong> {event.AdditionalInfo}
            </Card.Text>
          )}
        </div>

        {/* Правая часть: кнопки */}
        <div className="d-flex flex-column align-items-end">
          {role === UserRoles.Organaizer ? (
            <Button
              variant="outline-danger"
              className="mb-2"
              onClick={deleteEvent}
              style={{ display: "flex", alignItems: "center" }}
            >
              {/* <FaTrash className="me-2" /> */}
              Удалить
            </Button>

          ) : role === UserRoles.Student && isSubmitApplication && (
            <Button
              variant="outline-success"
              className="mb-2"
              onClick={onSubmitApplication}
              style={{ display: "flex", alignItems: "center" }}
            >
              Подать заявку
            </Button>
          )}
          {/* Кнопка для перехода на страницу заявок */}
          {role === UserRoles.Organaizer && isSubmitApplication && (
            <Link to={`/applications/event/${event.ID}`}>
              <Button
                variant="outline-primary"
                className="mb-2"
                style={{ display: "flex", alignItems: "center" }}
              >
                Перейти к заявкам
              </Button>
            </Link>
          )}

        </div>
      </Card.Body>
      {event.Events && event.Events.length > 0 && (
        <Card.Footer>
          <div>
            <button
              className="btn btn-primary"
              type="button"
              data-bs-toggle="collapse"
              data-bs-target={`#collapse-${event.ID}`}
              aria-expanded="false"
              aria-controls={`collapse-${event.ID}`}
            >
              Показать вложенные события
            </button>
            <div className="collapse" id={`collapse-${event.ID}`}>
              <div className="d-flex">
                {event.Events.map((childEvent) => (
                  <Card key={childEvent.ID} className="col m-1 mb-3">
                    <Card.Body>
                      <Card.Title>{childEvent.Name}</Card.Title>
                      <Card.Text>
                        {new Date(childEvent.StartDate).toLocaleString()} -  {new Date(childEvent.EndDate).toLocaleString()}
                      </Card.Text>
                      {childEvent.AdditionalInfo && (
                        <Card.Text>
                          Дополнительная информация: {childEvent.AdditionalInfo}
                        </Card.Text>
                      )}
                    </Card.Body>
                  </Card>
                ))}
              </div>
            </div>
          </div>
        </Card.Footer>
      )}
    </Card>
  )
}

export default EventItem