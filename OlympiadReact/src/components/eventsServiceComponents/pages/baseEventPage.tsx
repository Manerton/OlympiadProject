import EventList from "../eventList"

import { useState, useEffect, useMemo } from "react";
import { Button, Modal } from "react-bootstrap";
import EventModalForm from "../eventModalWindow";
import { CLASS, MyEvent, OLYMPIAD, REGIONAL_STAGE, STAGE } from "../../../types/event";
import { useRole } from "../../RoleContext";
import API_CONFIG from "../../../config/apiConfig"
import EventInfo from "../eventInfo";
import UserRoles from "../../../types/user";
import Pagination from "../pagination";
// import { useUser } from "../contexts/UserContext";

interface BaseEventPageProps {
  selectedEventId?: string
  pageName: string
  showSubjectField?: boolean
  showClassNumber?: boolean
  EventType: string
}

const limit = 10


function BaseEventPage({ selectedEventId, pageName, EventType, showSubjectField = false, showClassNumber = false }: BaseEventPageProps) {
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [showModal, setShowModal] = useState(false);

  const [events, setEvents] = useState<MyEvent[]>([])
  const [event, setEvent] = useState<MyEvent>()

  const { role } = useRole();
  const [eventFieldName, setEventFieldName] = useState("")
  

  const memoizedShowSubjectField = useMemo(() => showSubjectField, [showSubjectField]);
  const memoizedShowClassNumberField = useMemo(() => showClassNumber, [showClassNumber]);



  const eventsResponse = async (offset: number) => {
    let endPointEvents = ""
    switch (EventType) {
      case REGIONAL_STAGE:
        endPointEvents = `${API_CONFIG.EVENTS}/regional-stage?offset=${offset}&limit=${limit}`;
        break
      case OLYMPIAD:
        endPointEvents = `${API_CONFIG.EVENTS}/child/${selectedEventId}?offset=${offset}&limit=${limit}`;
        break
      case CLASS:
        endPointEvents = `${API_CONFIG.EVENTS}/child/${selectedEventId}?offset=${offset}&limit=${limit}`;
        break
      case STAGE:
        endPointEvents = `${API_CONFIG.EVENTS}/stages/${selectedEventId}`;
    }
    try {
      const response = await fetch(endPointEvents, {
        method: "GET",
        credentials: "include", // Отправка cookie
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const result = await response.json();
      console.log("EVENTS: Response from API:", result);
      if (EventType === STAGE && result.data) {
        const eventsWithDates = result.data.map((event: MyEvent) => ({
          ...event,
          StartDate: new Date(event.start_date),
          EndDate: new Date(event.end_date),
        }));
        setEvents(eventsWithDates);
        setTotalPages(0)
      } else {
        const eventsWithDates = result.data.events.map((event: MyEvent) => ({
          ...event,
          StartDate: new Date(event.start_date),
          EndDate: new Date(event.end_date),
        }));
        setEvents(eventsWithDates);
        const totalCount = result.data.totalCount
        const totalPage = totalCount / limit
        setTotalPages(totalCount%limit == 0 ? totalPage : totalPage + 1)
      }
    } catch (error) {
      console.error("Ошибка при загрузке региональных этапов:", error);
    }
  }

  const eventResponse = async () => {
    const endPointEvent = `${API_CONFIG.EVENTS}/${selectedEventId}`;
    if (selectedEventId) {
      try {
        const response = await fetch(endPointEvent, {
          method: "GET",
          credentials: "include", // Отправка cookie
        });
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const result = await response.json();
        console.log("EVENT: Response from API:", result);

        if (result.data) {
          const eventWithDates = {
            ...result.data,
            StartDate: new Date(result.data.StartDate),
            EndDate: new Date(result.data.EndDate),
          };
          setEvent(eventWithDates);

          console.log(eventWithDates)

        }
      } catch (error) {
        console.error("Ошибка при загрузке региональных этапов:", error);
      }
    }
  }

  const updateBySetPage = async(newPage: number) => {
    setPage(newPage)
    let offset = newPage * limit - limit
    eventsResponse(offset)
  }


  const fetchEvents = async () => {
    await eventsResponse(0)
    await eventResponse()
  };

  useEffect(() => {
    switch (EventType) {
      case OLYMPIAD: 
        setEventFieldName("олимпиады")
        break 
      case STAGE:
        setEventFieldName("этапа олимпиады")
        break
      case STAGE:
        setEventFieldName("стадии этапа")
        break
      case REGIONAL_STAGE:
        setEventFieldName("регионального этапа")
        break
      default:
        setEventFieldName("события")
    }

    fetchEvents();
  }, [selectedEventId]);

  const OnUpdateListEvent = () => {
    setShowModal(false)
    fetchEvents()
  }

  const handleDeleteEvent = (id: number) => {
    setEvents((events) => events.filter((event) => event.id !== id));
  };

  const sortEvents = (order: "asc" | "desc") => {
    setEvents((prevEvents) =>
      [...prevEvents].sort((a, b) => {
        const dateA = a.start_date.getTime();
        const dateB = b.start_date.getTime();
        return order === "asc" ? dateA - dateB : dateB - dateA;
      })
    );
  };

  return (
    <div className="row d-flex justify-content-center">
      {event && (
        <div className="col-3">
          <h3>Информация о {event.name}</h3>
          <EventInfo event={event}></EventInfo>
        </div>
      )}
      <div className="col-9">
        <div className="d-flex justify-content-between align-items-center">
          <h1>{pageName}</h1>
          <div className="d-flex">
            <div className="dropdown">
              <button
                className="btn btn-secondary dropdown-toggle"
                type="button"
                data-bs-toggle="dropdown"
                aria-expanded="false"
              >
                Сортировать
              </button>
              <ul className="dropdown-menu">
                <li>
                  <a
                    className="dropdown-item"
                    onClick={() => sortEvents("asc")}
                    style={{ cursor: "pointer" }}
                  >
                    По дате возрастания
                  </a>
                </li>
                <li>
                  <a
                    className="dropdown-item"
                    onClick={() => sortEvents("desc")}
                    style={{ cursor: "pointer" }}
                  >
                    По дате убывания
                  </a>
                </li>
              </ul>
            </div>


            {/* Кнопка создания доступна только организаторам */}
            {role === UserRoles.Organaizer && EventType != OLYMPIAD && (
              <Button
                variant="primary"
                className="ms-2"
                onClick={() => setShowModal(true)}
              >
                Создать
              </Button>
            )}
          </div>
        </div>

        {/* Список этапов */}
        {events.length > 0 ? (
          <EventList isSubmitItem={EventType === OLYMPIAD ? true : false} events={events} onDelete={handleDeleteEvent} />
        ) : (
          <p>Список пуст...</p>
        )}
        {/* Пагинация */}
        <Pagination currentPage={page} totalPages={totalPages} onPageChange={updateBySetPage} />
        {/* Модальное окно */}
      </div>
      <Modal show={showModal} onHide={() => setShowModal(false)}>
        <Modal.Header closeButton>
          <Modal.Title>Создание {eventFieldName}</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          {/* Форма создания */}
          <EventModalForm onSuccess={OnUpdateListEvent} event={event} showSubjectField={memoizedShowSubjectField} showClassNumber={memoizedShowClassNumberField} ></EventModalForm>
        </Modal.Body>
      </Modal>
    </div>
  );
}

export default BaseEventPage;
