import EventList from "../eventList"

import React, { useState, useEffect, useMemo } from "react";
import { Button, Modal } from "react-bootstrap";
import EventModalForm from "../eventModalWindow";
import { MyEvent, OLYMPIAD, REGIONAL_STAGE, STAGE } from "../../../types/event";
import { RoleProvider, useRole } from "../../RoleContext";
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

  const [events, setEvents] = useState([])
  const [event, setEvent] = useState<MyEvent>()
  const [isLoading, setIsLoading] = useState(true);

  const { role, id } = useRole();

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
        setEvents(result.data);
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
        setEvent(result.data);

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

  return (

    <div className="row">

      {event && (
        <div className="col-3">
          <h3>Информация о {event.Name}</h3>
          <EventInfo event={event}></EventInfo>
        </div>
      )}
        <div className="col-9">
          <div className="d-flex justify-content-between">
            <h1>{pageName}</h1>
            {/* Кнопка создания доступна только организаторам */}
            <p>{role}</p>
            {role === "3" && (
              <Button
                variant="primary"
                className="mb-3"
                onClick={() => setShowModal(true)}
              >
                Создать
              </Button>
            )}
          </div>


          {/* Список этапов */}
          {events.length > 0 ? (
            <EventList events={events} />
          ) : (
            <p>Загрузка данных...</p>
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
