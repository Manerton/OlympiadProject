import EventList from "../eventList"

import React, { useState, useEffect } from "react";
import { Button, Modal } from "react-bootstrap";
import EventModalForm from "../eventModalWindow";
import { MyEvent, OLYMPIAD, REGIONAL_STAGE, STAGE } from "../../../types/event";
import { RoleProvider, useRole } from "../../RoleContext";
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
  const [event, setEvent] = useState()

  const { role, id } = useRole();

  const fetchEvents = async () => {
    let endPointEvents = ""
    if (type === REGIONAL_STAGE) {
      endPointEvents = `http://localhost:8080/events/regional-stage`;
    } else if (type === OLYMPIAD) {
      endPointEvents = `http://localhost:8080/events/child/${selectedEventId}`
    } else if (type === STAGE) {
      endPointEvents = `http://localhost:8080/events/stages/${selectedEventId}`;
    }


    const endPointEvent = `http://localhost:8080/events/${selectedEventId}`;
    try {

      const response = await fetch(endPointEvents, {
        credentials: "include", // Отправка cookie
    });
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const result = await response.json();
      console.log("Response from API:", result);
      setEvents(result.data);
    } catch (error) {
      console.error("Ошибка при загрузке региональных этапов:", error);
    }


    if (selectedEventId) {
      try {
        const response = await fetch(endPointEvent);
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const result = await response.json();
        console.log("Response from API:", result);
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
        
      <div>
        <div className="d-flex justify-content-between">
          <h1>{pageName}</h1>
          {/* Кнопка создания доступна только организаторам */}
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
        <EventList events={events} parentEvent={event} />
        {/* Пагинация */}
        {/* <Pagination currentPage={page} totalPages={totalPages} onPageChange={setPage} /> */}
        {/* Модальное окно */}
        <Modal show={showModal} onHide={() => setShowModal(false)}>
          <Modal.Header closeButton>
            <Modal.Title>Создать</Modal.Title>
          </Modal.Header>
          <Modal.Body>
            {/* Форма создания */}
            <EventModalForm onSuccess={OnUpdateListEvent} event={event} showSubjectField={showSubjectField} ></EventModalForm>
          </Modal.Body>
        </Modal>
      </div>
    );
}

export default BaseEventPage;
