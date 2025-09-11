// src/pages/StagesListPage.tsx
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { Container, Row, Col, Card, Button, Spinner, Alert } from "react-bootstrap";
import type { MyEvent } from "../../../types/event.ts";
import { fetchOlympiadStages } from "../../../../requests/EventsRequests";
import { useAuth } from "../../../Helpers/AuthContext";
import { UserRole } from "../../../../dictionary/role.js";
import { EventType } from "../../../../dictionary/eventType.js";

interface JuryMember {
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
  return futureEvents[0].id!
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

const StagesListPage: React.FC = () => {
  const { id } = useParams();
  const [stages, setStages] = useState<MyEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { user } = useAuth();

  useEffect(() => {
    if (!id) return;
    setLoading(true);
    fetchOlympiadStages(id)
      .then((res) => {
        setStages(res ?? []);
      })
      .catch((err) => setError((err as Error).message))
      .finally(() => setLoading(false));
  }, [id]);

  if (loading) return <Spinner animation="border" />;
  if (error) return <Alert variant="danger">{error}</Alert>;
  if (!stages || stages.length === 0) return <Alert variant="warning">Этапы не найдены</Alert>;

  return (
    <Container className="py-4">
      <h3 className="fw-bold mb-4">Этапы олимпиады</h3>
      <Row>
        {stages.map((stage) => {

          const nextEventId = findNextEventId(stage);
          // пока моки для жюри, можно будет заменить на stage.jury
          const jury: JuryMember[] = [
            { name: "Иван Иванов", role: "Председатель" },
            { name: "Мария Петрова", role: "Член жюри" },
            { name: "Сергей Кузнецов", role: "Член жюри" },
            { name: "Сергей Кузнецов", role: "Член жюри" },
          ];

          return (
            <Col md={12} key={stage.id} className="mb-4">
              <div className="shadow-sm mb-4">
                <div className="row">
                  <div className="col-8">
                    <div className="d-flex align-items-center gap-3 mt-3 flex-wrap">
                      {/* Этап */}
                      <div
                        className={`p-3 border rounded text-center card-stage ${nextEventId === stage.id ? "border-primary border-3 shadow" : ""
                          }`}
                      >
                        <h6 className="fw-bold mb-3">{stage.name}</h6>
                        <hr />
                        <div className="justify-content-between">
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
                            className={`p-3 border rounded text-center card-stage ${nextEventId === view.id ? "border-primary border-3 shadow" : ""
                              }`}
                          >
                            <h6 className="fw-bold mb-3">Просмотр работ</h6>
                            <hr />

                            <div className="justify-content-between">
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
                            className={`p-3 border rounded text-center card-stage ${nextEventId === appeal.id ? "border-primary border-3 shadow" : ""
                              }`}
                          >
                            <h6 className="fw-bold mb-3">Апелляция</h6>
                            <hr />

                            <div className="justify-content-between">
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

                  </div>
                  <div className="col-4">
                    <h6 className="mt-3 fw-bold">Список жюри</h6>
                    {jury.map((member, idx) => (
                      <p key={idx} className="mb-1">
                        <strong>{member.name}</strong> — {member.role}
                      </p>
                    ))}
                    {user?.role === UserRole.Admin && (
                      <Button variant="primary" className="w-100 mt-3">
                        Назначить жюри
                      </Button>
                    )}
                  </div>

                </div>

              </div>
              <hr />
            </Col>
          );
        })}
      </Row>
    </Container>
  );
};

export default StagesListPage;
