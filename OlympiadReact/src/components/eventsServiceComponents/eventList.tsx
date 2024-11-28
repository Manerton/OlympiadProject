import EventItem from "./eventItem";
import { MyEvent } from "../../types/event";

interface EventListProps {
  events: MyEvent[]
  parentEvent?: MyEvent
}



function EventList({events, parentEvent}: EventListProps) {
  


    if (!events) {
      return <div>
         {parentEvent && (
              <div>
                <p>{new Date(parentEvent.StartDate).toLocaleString()} -- {new Date(parentEvent.EndDate).toLocaleString()}</p>
              </div>
            )}
          Список пуст...
        </div>;
    }

    return (
        

        <div>
            {parentEvent && (
              <div>
                <p>{new Date(parentEvent.StartDate).toLocaleString()} -- {new Date(parentEvent.EndDate).toLocaleString()}</p>
              </div>
            )}
            {events.map((event: MyEvent) => (
                <EventItem key={event.ID} event={event} />
            ))}
        </div>
    )
}

export default EventList