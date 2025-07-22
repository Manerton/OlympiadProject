import EventItem from "./eventItem";
import { MyEvent } from "../../types/event";

interface EventListProps {
  events: MyEvent[]
  onDelete:  (id: number) => void;
  isSubmitItem?: boolean
}

function EventList({ events, onDelete, isSubmitItem}: EventListProps) {
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
        <EventItem isSubmitApplication={isSubmitItem || false} key={event.id} event={event} onDelete={onDelete} />
      ))}
    </div>
  )
}

export default EventList