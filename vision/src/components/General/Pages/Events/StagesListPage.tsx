// src/pages/OlympiadDetails.tsx
import React, { use, useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { Container, Row, Col, Card, Button, Spinner, Alert } from "react-bootstrap";
import type { MyEvent } from "../../../types/event.ts";
import CardImage from "./components/CardImage.tsx"; 
import { fetchEvent,fetchStagesCount} from "../../../../requests/EventsRequests.ts";

const StagesListPage: React.FC = () => {
  const { id } = useParams();
  const [olympiad, setOlympiad] = useState<MyEvent | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [stagesCount, setStagesCount] = useState<number | null>(null);

  useEffect(() => {
  if (!id) return;
    setLoading(true);
    fetchEvent(id)   // ✅ вот здесь
        .then((res) => setOlympiad(res))
        .catch((err) => setError((err as Error).message))
        .finally(() => setLoading(false));

        fetchStagesCount(id)
        .then((count) => setStagesCount(count))
        .catch((err) => console.error("Ошибка загрузки этапов:", err));

    }, [id]);

  if (loading) return <Spinner animation="border" />;
  if (error) return <Alert variant="danger">{error}</Alert>;
  if (!olympiad) return <Alert variant="warning">Олимпиада не найдена</Alert>;

  return (
    <Container className="py-4">
      <Row>
        {/* Левая часть */}
        <Col md={8} className="shadow-sm border-1">
          <CardImage
            subjectId={olympiad.subject ?? 0}
            width={800}
            height={300}
          />
          <h3 className="fw-bold mt-3">{olympiad.name}</h3>
          <h6 className="text-uppercase text-secondary fw-bold">Описание</h6>
          <p>{olympiad.additional_info}</p>
        </Col>

        {/* Правая карточка */}
        <Col md={4}>
          <Card className="shadow-sm">
            <CardImage
              subjectId={olympiad.subject ?? 0}
              width={400}
              height={150}
            />
            <Card.Body>
              <Button variant="primary" className="w-100 mb-3">
                Подать заявку
              </Button>
              <p><strong>Дата начала:</strong> {olympiad.start_date}</p>
              <p><strong>Время:</strong> {olympiad.end_date}</p>
               <strong>Количество этапов:</strong>{" "}
                {stagesCount !== null ? stagesCount : "Загрузка..."}
              <p><strong>Предмет:</strong> {olympiad.subject}</p>
              <p><strong>Место проведения:</strong>???</p>
            </Card.Body>
          </Card>
        </Col>
      </Row>
    </Container>
  );
};

export default StagesListPage;
