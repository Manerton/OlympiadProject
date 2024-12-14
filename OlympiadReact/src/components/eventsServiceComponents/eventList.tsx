import EventItem from "./eventItem";
import { MyEvent } from "../../types/event";

interface EventListProps {
  events: MyEvent[]
  onDelete:  (id: number) => void;
}

function EventList({ events, onDelete }: EventListProps) {
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
        <EventItem key={event.ID} event={event} onDelete={onDelete} />
      ))}
    </div>
  )
}

export default EventList