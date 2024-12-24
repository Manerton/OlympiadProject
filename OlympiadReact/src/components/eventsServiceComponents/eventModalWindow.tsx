import React, { useEffect, useState } from "react";
import { Button, Form } from "react-bootstrap";
import { MyEvent, OLYMPIAD, REGIONAL_STAGE, STAGE } from "../../types/event.ts";
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
  const [startDate, setStartDate] = useState<string>(
    event?.StartDate ? formatDateForInput(new Date(event.StartDate).toISOString()) : ""
  );
  const [endDate, setEndDate] = useState<string>(
    event?.EndDate ? formatDateForInput(new Date(event.EndDate).toISOString()) : ""
  );
  const [additionalInfo, setAdditionalInfo] = useState("");



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

    const eventData: MyEvent = {
      PreviousEventID: event?.ID,
      Name: eventName,
      StartDate: new Date(startDate),
      EndDate: new Date(endDate),
      EventType: event?.EventType || "",
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
      if (result.status === "Error") {
        alert(`Не удалось создать событие: ${result.error}`);
      }
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
      <Form.Group className="mb-3">
        <Form.Label>Название</Form.Label>
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
          <Form.Select
            value={subject}
            onChange={(e) => setSubject(e.target.value)}
          >
            <option value="">Выберите предмет</option>
            {subjectList.map((subject_str) => (
              <option key={subject_str} value={subject_str}>
                {subject_str}
              </option>
            ))}
          </Form.Select>
        </Form.Group>
      )}

      <Form.Group className="mb-3">
        <Form.Label>Дата начала</Form.Label>
        <Form.Control
          type="datetime-local"
          value={startDate}
          onChange={(e) => setStartDate(e.target.value)}
        />
      </Form.Group>

      <Form.Group className="mb-3">
        <Form.Label>Дата конца</Form.Label>
        <Form.Control
          type="datetime-local"
          value={endDate}
          onChange={(e) => setEndDate(e.target.value)}
        />
      </Form.Group>

      <Form.Group className="mb-3">
        <Form.Label>Дополнительная информация</Form.Label>
        <Form.Control
          as="textarea"
          rows={3}
          value={additionalInfo}
          onChange={(e) => setAdditionalInfo(e.target.value)}
          placeholder="Введите дополнительную информацию"
        />
      </Form.Group>

      <Button variant="primary" type="submit">
        Создать
      </Button>
    </Form>
  );
}

export default EventModalForm;
