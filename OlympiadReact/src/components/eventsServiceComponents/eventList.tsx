import EventItem from "./eventItem";
import { MyEvent } from "../../types/event";
import { Button, Form } from "react-bootstrap";
import { useEffect, useState } from "react";
import formatDateForInput from "../../support/support";
import API_CONFIG from "../../config/apiConfig";


interface EventListProps {
  events: MyEvent[]
}



function EventList({ events }: EventListProps) {


  if (!events) {
    return (
      <div>
        Список пуст...
      </div>
    )
  }

  return (
    <div>
      {events.map((event: MyEvent) => (
        <EventItem key={event.ID} event={event} />
      ))}
    </div>
  )


}

export default EventList