import { useEffect, useState } from "react";
import type { ApplicationEvent } from "../../../../../types/event";
import ApplicationEventCart from "./ApplicationEventCart";

const ApplicationEventTab: React.FC = () => {
    const [events, setEvents] = useState<ApplicationEvent[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        async function fetchApplicationEvents(userId: string) {
            try {
                const result = await fetch("");
                if (!result.ok) 
                    throw new Error("Ошибка при загрузке олимпиад");
                const data: ApplicationEvent[] = await result.json();
                setEvents(data);
            } catch (err) {
                console.error(err)
            } finally {
                setLoading(false)
            }
        }

        fetchApplicationEvents("")
    }, [])


    if (loading) return <div>Загрузка...</div>

    if (events.length === 0) return <div>Нет заявок на участие</div>

    return (
        <div>
            {events.map((event) => (
                <ApplicationEventCart key={event.id} event={event}/>
            ))}
        </div>
    );
};

export default ApplicationEventTab;