import { MyEvent, REGIONAL_STAGE, STAGE, OLYMPIAD } from "../../types/event";
import { Button, Card } from "react-bootstrap";
import { useNavigate } from "react-router-dom";
import { useRole } from "../RoleContext";
import API_CONFIG from "../../config/apiConfig";

interface EventItemProps {
  event: MyEvent;
}

function EventItem({ event }: EventItemProps) {
  const navigate = useNavigate();
  const { role, id } = useRole()

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
    const parentDeletedEventID = event.PreviousEventID
    try {

      const response = await fetch(endPointDeleteEvent, { method: 'DELETE' });
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const result = await response.json();
      console.log("Response from API:", result);
    } catch (error) {
      console.error("Ошибка при загрузке региональных этапов:", error);
    }

    console.log("deleted ID:", parentDeletedEventID)
    if (event.EventType === REGIONAL_STAGE) {
      navigate(`/events`);
    } else if (event.EventType === OLYMPIAD) {
      navigate(`/olympiads/${parentDeletedEventID}`);
    } else if (event.EventType === STAGE) {
      navigate(`/olympiad-stages/${parentDeletedEventID}`);
    }
  }

  return (
    <Card
      key={event.ID}
      className="mb-3"
    //Переход на вложенный eventList по id
    >
      <Card.Body className="d-flex justify-content-between">
        <div onClick={handleClick} style={{ cursor: "pointer" }}>
          <Card.Title>{event.Name}</Card.Title>

          <Card.Text>Дата начала: {new Date(event.StartDate).toLocaleString()}</Card.Text>
          <Card.Text>Дата конца: {new Date(event.EndDate).toLocaleString()}</Card.Text>
          {event.AdditionalInfo && (
            <div>
              <hr />
              <Card.Text>Дополнительная информация: {event.AdditionalInfo}</Card.Text>
            </div>
          )}
        </div>
        <div>
          {role === "3" && (
            <Button variant="danger" onClick={deleteEvent} >
              Удалить
            </Button>
          )}
        </div>


      </Card.Body>
      <Card.Footer>
        {event.Events && event.Events.length > 0 && (
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
              <div className="card card-body">
                {event.Events.map((childEvent) => (
                  <Card key={childEvent.ID} className="mb-3">
                    <Card.Body>
                      <Card.Title>{childEvent.Name}</Card.Title>
                      <Card.Text>
                        Дата начала: {new Date(childEvent.StartDate).toLocaleString()}
                      </Card.Text>
                      <Card.Text>
                        Дата конца: {new Date(childEvent.EndDate).toLocaleString()}
                      </Card.Text>
                      {childEvent.AdditionalInfo && (
                        <Card.Text>
                          Дополнительная информация: {childEvent.AdditionalInfo}
                        </Card.Text>
                      )}
                      {/* Дополнительные действия для вложенных событий */}
                      <Button
                        variant="primary"
                        onClick={() => navigate(`/events/${childEvent.ID}`)}
                      >
                        Перейти к событию
                      </Button>
                    </Card.Body>
                  </Card>
                ))}
              </div>
            </div>
          </div>
        )}
      </Card.Footer>

    </Card>
  )
}

export default EventItem