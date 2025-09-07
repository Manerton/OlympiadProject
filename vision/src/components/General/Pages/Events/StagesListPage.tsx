// src/pages/StagesListPage.tsx
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { Container, Row, Col, Card, Button, Spinner, Alert } from "react-bootstrap";
import type { MyEvent } from "../../../types/event.ts";
import { fetchOlympiadChildren } from "../../../../requests/EventsRequests";
import { useAuth } from "../../../Helpers/AuthContext";

interface JuryMember {
  name: string;
  role: string;
}

const StagesListPage: React.FC = () => {
  const { id } = useParams();
  const [stages, setStages] = useState<MyEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { user } = useAuth();
  useEffect(() => {
    if (!id) return;
    setLoading(true);
    fetchOlympiadChildren(id)
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
          // пока моки для жюри, можно будет заменить на stage.jury
          const jury: JuryMember[] = [
            { name: "Иван Иванов", role: "Председатель" },
            { name: "Мария Петрова", role: "Член жюри" },
            { name: "Сергей Кузнецов", role: "Член жюри" },
          ];

          return (
            <Col md={12} key={stage.id} className="mb-4">
              <Card className="shadow-sm">
                <Card.Body>
                  <h5 className="fw-bold">{stage.name}</h5>
                  <p className="fs-5">{stage.additional_info}</p>
                  <p className="text-muted">
                    {stage.start_date} — {stage.end_date}
                  </p>

                  <h6 className="mt-3 fw-bold">Жюри</h6>
                  {jury.map((member, idx) => (
                    <p key={idx} className="mb-1">
                      <strong>{member.name}</strong> — {member.role}
                    </p>
                  ))}
                  
                  {user?.role === 1 && (
                  <Button variant="primary" className="w-100 mt-3">
                    Назначить жюри
                  </Button>
                )}

                  
                </Card.Body>
              </Card>
            </Col>
          );
        })}
      </Row>
    </Container>
  );
};

export default StagesListPage;
