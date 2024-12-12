import React, { useEffect, useState } from "react";
import { Button, Form } from "react-bootstrap";
import { MyEvent } from "../../types/event.ts";
import API_CONFIG from "../../config/apiConfig.ts";
import formatDateForInput from "../../support/support.ts";

export interface EventModalFormProps {
  showSubjectField?: boolean;
  event?: MyEvent;
  onSuccess?: () => void;
}

function EventModalForm({ event, showSubjectField, onSuccess }: EventModalFormProps) {
  const [eventName, setEventName] = useState("");
  const [subjectList, setSubjectList] = useState<string[]>([]);
  const [subject, setSubject] = useState("");
  const [startDate, setStartDate] = useState<string>(formatDateForInput(event?.StartDate || ""));
  const [endDate, setEndDate] = useState<string>(formatDateForInput(event?.EndDate || ""));
  const [additionalInfo, setAdditionalInfo] = useState<string>("");
  // const [subEventType, setSubEventType] = useState<string>("")

  useEffect(() => {
    const getSubjects = async () => {
      try {
        const response = await fetch(`${API_CONFIG.EVENTS}/subjects`, {
          method: "GET",
          credentials: "include", // Для отправки cookie
          headers: {
            "Content-Type": "application/json",
          },
        });

        if (!response.ok) {
          const errorText = await response.text();
          throw new Error(`Ошибка API: ${errorText}`);
        }

        const result = await response.json();
        console.log("Предметы получены!", result);
        setSubjectList(result.data);
      } catch (error) {
        console.error("Ошибка при получении предметов:", error);
      }
    };

    if (showSubjectField && subjectList.length === 0) {
      getSubjects();
    }
  }, [showSubjectField, subjectList.length]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!eventName || !startDate || !endDate) {
      alert("Заполните все поля!");
      return;
    }

    if (showSubjectField && !subject) {
      alert("Заполните поле Предмет!");
      return;
    }

    // if (showSelectSubEventType && !subEventType) {
    //   alert("Заполните поле Тип события!");
    //   return;
    // }

    const eventData: MyEvent = {
      PreviousEventID: event?.ID,
      Name: eventName,
      StartDate: new Date(startDate).toISOString(),
      EndDate: new Date(endDate).toISOString(),
      EventType: "",
      Subject: subject,
      AdditionalInfo: additionalInfo,
    };

    try {
      const response = await fetch(`${API_CONFIG.EVENTS}`, {
        method: "POST",
        credentials: "include", // Для отправки cookie
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(eventData),
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(`Ошибка API: ${errorText}`);
      }

      const result = await response.json();
      console.log("Событие создано!", result.Date);

      setEventName("");
      if (onSuccess) {
        onSuccess();
      }
    } catch (error) {
      console.error("Ошибка при создании события:", error);
      alert("Не удалось создать событие. Попробуйте снова.");
    }
  };

  return (
    <Form onSubmit={handleSubmit}>

      {/* {showSelectSubEventType && (
      <Form.Group className="mb-3">
        <Form.Label>Тип события</Form.Label>
        <Form.Control
          as="select"
          placeholder="Выберите тип"
          value={subEventType}
          onChange={(e) => setSubEventType(e.target.value)}>

          <option value="">Выберите тип</option>
          <option value="Просмотр работ">Просмотр работ</option>
          <option value="Апелляция">Апелляция</option>
        </Form.Control>

      </Form.Group>
      )} */}

      <Form.Group className="mb-3">
        <Form.Label>Название события</Form.Label>
        <Form.Control
          type="text"
          placeholder="Введите название"
          value={eventName}
          onChange={(e) => setEventName(e.target.value)}
        />
      </Form.Group>

      {showSubjectField && (
        <Form.Group controlId="eventSubject" className="mb-3">
          <Form.Label>Предмет</Form.Label>
          <div className="dropdown">
            <button
              className="btn btn-secondary dropdown-toggle"
              type="button"
              data-bs-toggle="dropdown"
              aria-expanded="false"
            >
              {subject || "Выберите предмет"}
            </button>
            <ul className="dropdown-menu">
              {subjectList.map((subject_str) => (
                <li>
                  <a
                    className="dropdown-item"
                    onClick={() => setSubject(subject_str)}
                    style={{ cursor: "pointer" }}
                  >
                    {subject_str}
                  </a>
                </li>
              ))}
            </ul>
          </div>
        </Form.Group>
      )}

      <Form.Group className="mb-3">
        <Form.Label>Дата начала события</Form.Label>
        <Form.Control
          type="datetime-local"
          value={startDate}
          min={ event?.EventType != "STAGE" ? formatDateForInput(event?.StartDate || "") : formatDateForInput(event?.EndDate || "")}
          max={ event?.EventType != "STAGE" ? formatDateForInput(event?.EndDate || "") : undefined }
          onChange={(e) => setStartDate(e.target.value)}
        />
      </Form.Group>

      <Form.Group className="mb-3">
        <Form.Label>Дата конца события</Form.Label>
        <Form.Control
          type="datetime-local"
          value={endDate}
          min={ event?.EventType != "STAGE" ? formatDateForInput(event?.StartDate || "") : formatDateForInput(event?.EndDate || "")}
          max={ event?.EventType != "STAGE" ? formatDateForInput(event?.EndDate || "") : undefined }
          onChange={(e) => setEndDate(e.target.value)}
        />
      </Form.Group>

      <Form.Group className="mb-3">
        <Form.Label>Дополнительная информация</Form.Label>
        <Form.Control
          as="textarea"
          rows={5}
          value={additionalInfo}
          onChange={(e) => setAdditionalInfo(e.target.value)}
          placeholder="Введите здесь дополнительную информацию"
        />
      </Form.Group>

      <Button variant="primary" type="submit">
        Создать
      </Button>
    </Form>
  );
}

export default EventModalForm;
