import { Button, Form } from "react-bootstrap"
import API_CONFIG from "../../config/apiConfig";
import { useEffect, useState } from "react";
import formatDateForInput from "../../support/support";
import { MyEvent } from "../../types/event";

interface EventInfoProps {
  event: MyEvent
}

function EventInfo({ event }: EventInfoProps) {
  const [eventName, setEventName] = useState<string>(event.Name || "");
  const [subjectList, setSubjectList] = useState<string[]>([]);
  const [subject, setSubject] = useState<string>(event.Subject || "");
  const [startDate, setStartDate] = useState<string>(formatDateForInput(event.StartDate || ""));
  const [endDate, setEndDate] = useState<string>(formatDateForInput(event.EndDate || ""));
  const [additionalInfo, setAdditionalInfo] = useState<string>(event.AdditionalInfo || "");

  useEffect(() => {
    if (event.Subject) {
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
      getSubjects();
    }
  }, [event.Subject]);


  const UpdateEvent = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!eventName || !startDate || !endDate) {
      alert("Заполните все поля!");
      return;
    }

    if (event.Subject && !subject) {
      alert("Заполните поле предмет")
      return
    }


    const eventData: MyEvent = {
      ID: event.ID,
      PreviousEventID: event.PreviousEventID,
      Name: eventName,
      StartDate: new Date(startDate).toISOString(),
      EndDate: new Date(endDate).toISOString(),
      EventType: event.EventType || "",
      Subject: subject,
      AdditionalInfo: additionalInfo,
    };

    try {
      const response = await fetch(`${API_CONFIG.EVENTS}/${event.ID}`, {
        method: "PUT",
        credentials: "include", // Для отправки cookie
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(eventData),
      })

      if (!response.ok) {
        const errorText = response.text();
        throw new Error(`Ошибка API: ${errorText}`);
      }
      const result = await response.json();
      console.log("Событие создано!", result.data);

    } catch (error) {
      console.error("Ошибка при создании события:", error);
      alert("Не удалось создать событие. Попробуйте снова.");
    }
  }
  return (
      <div className="card shadow-sm mb-3">
        <div className="card-body">
          <h5 className="card-title">Редактировать событие</h5>
          <Form onSubmit={UpdateEvent}>


            {/* Даты и предмет */}
            {/* Название события */}
            <div className="row">
              <Form.Group className="col-md-4 mb-2">
                <Form.Label>Название события</Form.Label>
                <Form.Control
                  type="text"
                  placeholder="Введите название"
                  value={eventName}
                  onChange={(e) => setEventName(e.target.value)}
                />
              </Form.Group>

              <Form.Group className="col-md-4 mb-2">
                <Form.Label>Дата начала</Form.Label>
                <Form.Control
                  type="datetime-local"
                  value={startDate}
                  onChange={(e) => setStartDate(e.target.value)}
                />
              </Form.Group>
              <Form.Group className="col-md-4 mb-2">
                <Form.Label>Дата конца</Form.Label>
                <Form.Control
                  type="datetime-local"
                  value={endDate}
                  onChange={(e) => setEndDate(e.target.value)}
                />
              </Form.Group>
            </div>
            {subject && (

              <Form.Group className="mb-3">
                <Form.Label>Предмет</Form.Label>
                <div className="dropdown">
                  <button
                    className="btn btn-secondary dropdown-toggle w-100"
                    type="button"
                    data-bs-toggle="dropdown"
                    aria-expanded="false"
                  >
                    {subject || "Выберите предмет"}
                  </button>
                  <ul className="dropdown-menu w-100">
                    {subjectList.map((subject_str, index) => (
                      <li key={index}>
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
              </Form.Group>)}

            {/* Дополнительная информация */}
            <Form.Group className="mb-3">
              <Form.Label>Дополнительная информация</Form.Label>
              <Form.Control
                as="textarea"
                rows={3}
                value={additionalInfo}
                onChange={(e) => setAdditionalInfo(e.target.value)}
              />
            </Form.Group>

            <div className="text-end">
              <Button variant="primary" type="submit">
                Сохранить
              </Button>
            </div>
          </Form>
        </div>
      </div>



  )
}

export default EventInfo