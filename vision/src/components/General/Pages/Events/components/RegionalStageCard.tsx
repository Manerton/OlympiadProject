// components/RegionalStageCard.tsx
import React from "react";
import { Card, Button } from "react-bootstrap";
interface MyEvent {
    id?: string;
    name: string;
    start_date: string;
    end_date: string;
    previous_event_id?: string;
    event_type: string;
    subject?: number;
    class_number?: number;
    additional_info?: string;
    events?: MyEvent[];
}

interface Props {
  stage: MyEvent;
  onClick?: (id: string) => void;
}

const RegionalStageCard: React.FC<Props> = ({ stage, onClick }) => {
  const startDate = new Date(stage.start_date).toLocaleDateString();
  const endDate = new Date(stage.end_date).toLocaleDateString();

  return (
    <Card className="mb-3 shadow-sm">
      <Card.Body>
        <Card.Title>{stage.name}</Card.Title>
        <Card.Text>
          <strong>Дата начала:</strong> {startDate} <br />
          <strong>Дата окончания:</strong> {endDate} <br />
          {stage.additional_info && (
            <>
              <strong>Доп. информация:</strong> {stage.additional_info}
            </>
          )}
        </Card.Text>
        {onClick && (
          <Button variant="primary" onClick={() => onClick(stage.id!)}>
            Открыть
          </Button>
        )}
      </Card.Body>
    </Card>
  );
};

export default RegionalStageCard;
