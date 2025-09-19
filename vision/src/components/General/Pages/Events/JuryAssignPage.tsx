// src/pages/JuryAssignPage.tsx
import { useEffect, useState } from "react";
import { useParams, useLocation, useNavigate } from "react-router-dom";
import { Container, Row, Col, Button, ListGroup, Spinner, Alert } from "react-bootstrap";
import type { MyEvent } from "../../../types/event.ts";
import { fetchAllJury } from "../../../../requests/EventsRequests";
import { useAuth } from "../../../Helpers/AuthContext";
import { EventType } from "../../../../dictionary/eventType.js";

interface JuryMember {
  id?: string;
  name: string;
  role: string;
}

// helper: находит ближайшее событие в будущем
const findNextEventId = (stage: MyEvent): string | null => {
  const now = new Date();

  // Собираем этап + его дочерние события
  const allEvents = [stage, ...(stage.events ?? [])];

  // Парсим даты и фильтруем только будущие
  const futureEvents = allEvents
    .map((e) => ({
      ...e,
      start: new Date(e.start_date),
    }))
    .filter((e) => e.start >= now);

  if (futureEvents.length === 0) return null;

  // Берём ближайшее
  futureEvents.sort((a, b) => a.start.getTime() - b.start.getTime());
  return futureEvents[0].id!;
};

const formatDateTime = (dateString: string): string => {
  const date = new Date(dateString);
  return date.toLocaleDateString("ru-RU", {
    day: "2-digit",
    month: "long",
    hour: "2-digit",
    minute: "2-digit",
  });
};

const JuryAssignPage: React.FC = () => {
  const { id: stageId } = useParams<{ id: string }>();
  const location = useLocation();
  const navigate = useNavigate();
  const { user,accessToken } = useAuth();

  const [stage, setStage] = useState<MyEvent | null>(null);
  const [assigned, setAssigned] = useState<JuryMember[]>([]);
  const [available, setAvailable] = useState<JuryMember[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const state = location.state as { stage: MyEvent; olympiadId: string; assignedJury: JuryMember[] } | null;
    if (state) {
      setStage(state.stage);
      setAssigned(state.assignedJury);
      fetchAllJury(accessToken!, "2")//TODO ДОБАВИТЬ СЛОВАРИК ДЛЯ РОЛИ
        .then(setAvailable)
        .catch((err) => setError((err as Error).message))
        .finally(() => setLoading(false));
    } else {
      setError("Данные не переданы");
      setLoading(false);
    }
  }, [location.state]);

  if (loading) return <Spinner animation="border" />;
  if (error) return <Alert variant="danger">{error}</Alert>;
  if (!stage) return <Alert variant="warning">Данные этапа не найдены</Alert>;

  const nextEventId = findNextEventId(stage);

  const availableFiltered = available.filter(
    (m) => !assigned.some((a) => a.id === m.id)
  );

  const addToAssigned = (member: JuryMember) => {
    setAssigned([...assigned, { ...member, role: member.role || "Член жюри" }]);
  };

  const removeFromAssigned = (member: JuryMember) => {
    setAssigned(assigned.filter((a) => a.id !== member.id));
  };

  const onSave = () => {
    // Здесь можно добавить запрос на сохранение, например:
    // assignJuryToStage(stageId!, assigned)
    //   .then(() => navigate(-1))
    //   .catch((err) => setError(err.message));
    // Для плейсхолдера:
    alert("Изменения сохранены (плейсхолдер)");
    navigate(-1);
  };

  return (
    <Container className="py-4">
      <h3 className="fw-bold mb-4">Назначение жюри для этапа {stage.name}</h3>
      <div className="shadow-sm mb-4">
        <Row>
          <Col md={12}>
            <div className="d-flex align-items-center gap-3 mt-3 flex-wrap">
              {/* Этап */}
              <div
                className={`p-3 border rounded text-center card-stage ${
                  nextEventId === stage.id ? "border-primary border-3 shadow" : ""
                }`}
              >
                <h6 className="fw-bold mb-3">{stage.name}</h6>
                <hr />
                <div className="d-flex justify-content-between">
                  <div>
                    <span className="badge bg-success mb-1">Начало</span>
                    <p className="mb-0 fw-semibold">{formatDateTime(stage.start_date)}</p>
                  </div>
                  <div>
                    <span className="badge bg-danger mb-1">Конец</span>
                    <p className="mb-0 fw-semibold">{formatDateTime(stage.end_date)}</p>
                  </div>
                </div>
              </div>

              {/* ➝ если есть просмотр */}
              {stage.events?.some((c) => c.event_type === EventType.ViewWorks) && (
                <span className="fs-3">➝</span>
              )}

              {/* Просмотр работ */}
              {stage.events
                ?.filter((c) => c.event_type === EventType.ViewWorks)
                .map((view) => (
                  <div
                    key={view.id}
                    className={`p-3 border rounded text-center card-stage ${
                      nextEventId === view.id ? "border-primary border-3 shadow" : ""
                    }`}
                  >
                    <h6 className="fw-bold mb-3">Просмотр работ</h6>
                    <hr />

                    <div className="d-flex justify-content-between">
                      <div>
                        <span className="badge bg-success mb-1">Начало</span>
                        <p className="mb-0 fw-semibold">{formatDateTime(view.start_date)}</p>
                      </div>
                      <div>
                        <span className="badge bg-danger mb-1">Конец</span>
                        <p className="mb-0 fw-semibold">{formatDateTime(view.end_date)}</p>
                      </div>
                    </div>
                  </div>
                ))}

              {/* ➝ если есть апелляция */}
              {stage.events?.some((c) => c.event_type === EventType.Appeal) && (
                <span className="fs-3">➝</span>
              )}

              {/* Апелляция */}
              {stage.events
                ?.filter((c) => c.event_type === EventType.Appeal)
                .map((appeal) => (
                  <div
                    key={appeal.id}
                    className={`p-3 border rounded text-center card-stage ${
                      nextEventId === appeal.id ? "border-primary border-3 shadow" : ""
                    }`}
                  >
                    <h6 className="fw-bold mb-3">Апелляция</h6>
                    <hr />

                    <div className="d-flex justify-content-between">
                      <div>
                        <span className="badge bg-success mb-1">Начало</span>
                        <p className="mb-0 fw-semibold">{formatDateTime(appeal.start_date)}</p>
                      </div>
                      <div>
                        <span className="badge bg-danger mb-1">Конец</span>
                        <p className="mb-0 fw-semibold">{formatDateTime(appeal.end_date)}</p>
                      </div>
                    </div>
                  </div>
                ))}
            </div>
          </Col>
        </Row>
      </div>
      <hr />
      <Row>
        <Col md={6}>
          <h5>Назначенные жюри</h5>
          <ListGroup>
            {assigned.map((member, idx) => (
              <ListGroup.Item key={member.id || idx} className="d-flex justify-content-between align-items-center">
                <span>
                  <strong>{member.name}</strong> — {member.role}
                </span>
                <Button variant="danger" size="sm" onClick={() => removeFromAssigned(member)}>
                  Удалить
                </Button>
              </ListGroup.Item>
            ))}
            {assigned.length === 0 && <ListGroup.Item>Жюри не назначено</ListGroup.Item>}
          </ListGroup>
        </Col>
        <Col md={6}>
          <h5>Доступные жюри</h5>
          <ListGroup>
            {availableFiltered.map((member, idx) => (
              <ListGroup.Item key={member.id || idx} className="d-flex justify-content-between align-items-center">
                <span>
                  <strong>{member.name}</strong> {member.role ? `— ${member.role}` : ""}
                </span>
                <Button variant="success" size="sm" onClick={() => addToAssigned(member)}>
                  Добавить
                </Button>
              </ListGroup.Item>
            ))}
            {availableFiltered.length === 0 && <ListGroup.Item>Нет доступных</ListGroup.Item>}
          </ListGroup>
        </Col>
      </Row>
      <Button variant="primary" className="mt-4" onClick={onSave}>
        Сохранить изменения
      </Button>
      <Button variant="secondary" className="mt-4 ms-2" onClick={() => navigate(-1)}>
        Отмена
      </Button>
    </Container>
  );
};

export default JuryAssignPage;