
import React from 'react';
import { Card, Button } from 'react-bootstrap';
import type { MyEvent } from '../type/Event.ts'; // Assuming Event.ts is in the same directory

const OlympiadCard: React.FC<{ olympiad: MyEvent; onClick?: (id: string) => void }> = ({ olympiad, onClick }) => {
  const startDate = new Date(olympiad.start_date).toLocaleDateString();
  const endDate = new Date(olympiad.end_date).toLocaleDateString();

  return (
    <Card className="mb-3 shadow-sm">
      <Card.Img variant="top" src={`https://via.placeholder.com/300x200?text=${encodeURIComponent(olympiad.name)}`} />
      <Card.Body>
        <Card.Title>{olympiad.name}</Card.Title>
        <Card.Text>
          <strong>Дата начала:</strong> {startDate} <br />
          <strong>Дата окончания:</strong> {endDate} <br />
          <strong>Класс:</strong> {olympiad.class_number || 'N/A'} <br />
          <strong>Предмет:</strong> {olympiad.subject || 'N/A'} <br />
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