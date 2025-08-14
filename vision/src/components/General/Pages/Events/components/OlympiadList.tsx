import React from 'react';
import { Row, Col } from 'react-bootstrap';
import OlympiadCard from './OlympiadCard';
import { MyEvent } from '../type/Event.ts'; // Assuming Event.ts is in the same directory

const OlympiadList: React.FC<{ events: MyEvent[]; onEventClick?: (event: MyEvent) => void }> = ({ events, onEventClick }) => (
  <Row>
    {events.map((event) => (
      <Col md={4} key={event.id}>
        <EventCard 
          event={event} 
          onClick={event.event_type === 'REGIONAL_STAGE' && onEventClick ? () => onEventClick(event) : undefined} 
        />
      </Col>
    ))}
  </Row>
);

export default OlympiadList;