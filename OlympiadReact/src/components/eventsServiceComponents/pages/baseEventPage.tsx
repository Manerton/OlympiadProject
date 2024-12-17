import EventList from "../eventList"

import React, { useState, useEffect, useMemo } from "react";
import { Button, Modal } from "react-bootstrap";
import EventModalForm from "../eventModalWindow";
import { MyEvent, OLYMPIAD, REGIONAL_STAGE, STAGE } from "../../../types/event";
import { useRole } from "../../RoleContext";
import API_CONFIG from "../../../config/apiConfig"
import EventInfo from "../eventInfo";
// import Pagination from "./components/Pagination";
// import { useUser } from "../contexts/UserContext";

interface BaseEventPageProps {
  selectedEventId?: string
  pageName: string
  showSubjectField?: boolean
  type: string
}


function BaseEventPage({ selectedEventId, pageName, type, showSubjectField = false }: BaseEventPageProps) {
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [showModal, setShowModal] = useState(false);
  //   const { user } = useUser(); // Доступ к информации о пользователе

  const [events, setEvents] = useState<MyEvent[]>([])
  const [event, setEvent] = useState<MyEvent>()

  const { role } = useRole();

  const memoizedShowSubjectField = useMemo(() => showSubjectField, [showSubjectField]);

  const fetchEvents = async () => {
    let endPointEvents = ""
    switch (type) {
      case REGIONAL_STAGE:
        endPointEvents = `${API_CONFIG.EVENTS}/regional-stage`;
        break
      case OLYMPIAD:
        endPointEvents = `${API_CONFIG.EVENTS}/child/${selectedEventId}`;
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
      if (result.data) {
        const eventsWithDates = result.data.map((event: MyEvent) => ({
          ...event,
          StartDate: new Date(event.StartDate),
          EndDate: new Date(event.EndDate),
        }));
        setEvents(eventsWithDates);
      }
    } catch (error) {
      console.error("Ошибка при загрузке региональных этапов:", error);
    }

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
        }
      } catch (error) {
        console.error("Ошибка при загрузке региональных этапов:", error);
      }
    }
  };

  useEffect(() => {
    fetchEvents();
  }, [selectedEventId]);

  const OnUpdateListEvent = () => {
    setShowModal(false)
    fetchEvents()
  }

  const handleDeleteEvent = (id: number) => {
    setEvents((events) => events.filter((event) => event.ID !== id));
  };

  const sortEvents = (order: "asc" | "desc") => {
    setEvents((prevEvents) =>
      [...prevEvents].sort((a, b) => {
        const dateA = a.StartDate.getTime();
        const dateB = b.StartDate.getTime();
        return order === "asc" ? dateA - dateB : dateB - dateA;
      })
    );
  };


  return (
    <div className="row d-flex justify-content-center">
      {event && (
        <div className="col-3">
          <h3>Информация о {event.Name}</h3>
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
            {role === "3" && (
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
          <EventList isSubmitItem={type === OLYMPIAD ? true : false} events={events} onDelete={handleDeleteEvent} />
        ) : (
          <p>Список пуст...</p>
        )}
        {/* Пагинация */}
        {/* <Pagination currentPage={page} totalPages={totalPages} onPageChange={setPage} /> */}
        {/* Модальное окно */}
      </div>
      <Modal show={showModal} onHide={() => setShowModal(false)}>
        <Modal.Header closeButton>
          <Modal.Title>Создать</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          {/* Форма создания */}
          <EventModalForm onSuccess={OnUpdateListEvent} event={event} showSubjectField={memoizedShowSubjectField} ></EventModalForm>
        </Modal.Body>
      </Modal>
    </div>
  );
}

export default BaseEventPage;
