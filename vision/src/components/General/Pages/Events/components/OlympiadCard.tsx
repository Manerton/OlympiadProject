import React from 'react';
import { Card, Button } from 'react-bootstrap';
import type  { MyEvent } from '../../../../types/event.ts';  
import CardImage  from './CardImage.tsx'
const OlympiadCard: React.FC<{ olympiad: MyEvent; onClick?: (id: string) => void }> = ({ olympiad, onClick }) => {
  const startDate = new Date(olympiad.start_date).toLocaleDateString();
  const endDate = new Date(olympiad.end_date).toLocaleDateString();

  return (
    <Card className="mb-3 shadow-sm">
      <Card>
      <CardImage subjectId={olympiad.subject ?? 0} width={300} height={200} />
      <Card.Body>
        <Card.Title>{olympiad.name}</Card.Title>
      </Card.Body>
    </Card>
      <Card.Body>
        <Card.Text>
          <strong>Дата начала:</strong> {startDate} <br />
          <strong>Дата окончания:</strong> {endDate} <br />
          <strong>Класс:</strong> {olympiad.class_number || 'N/A'} <br />
          {olympiad.additional_info && (
            <>
              <strong>Доп. информация:</strong> {olympiad.additional_info} <br />
            </>
          )}
        </Card.Text>
        {onClick && (
          <Button variant="primary" onClick={() => onClick(olympiad.id!)}>
            Открыть
          </Button>
        )}
      </Card.Body>
    </Card>
  );
};

export default OlympiadCard;