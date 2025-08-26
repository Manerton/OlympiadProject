import { useEffect, useState, type ReactNode } from "react";
import type { ApplicationEvent } from "../../../../../types/event";
import ApplicationEventCart from "../EventCart";
import { StatusIcon } from "../../../../../Helpers/StatusBlock";
import { ApplicationStatus } from "../../../../../../dictionary/applicationStatus";
import { useNavigate } from "react-router-dom";

const AppealTab:  React.FC = () => {
    const [events, setEvents] = useState<ApplicationEvent[]>([]);
    const [loading, setLoading] = useState(true);

    const navigate = useNavigate()

    const mockEvents: ApplicationEvent[] = [
        {
            id: "1",
            name: "Региональный этап по математике",
            start_date: "2025-09-01T09:00:00Z",
            end_date: "2025-09-05T15:00:00Z",
            previous_event_id: "",
            subject: 1,
            class_number: 10,
            additional_info: "Проходит в онлайн-формате",
            status: ApplicationStatus.Pending,
        },
        {
            id: "2",
            name: "Олимпиада по физике",
            start_date: "2025-10-10T09:00:00Z",
            end_date: "2025-10-12T14:00:00Z",
            previous_event_id: "1",
            subject: 2,
            class_number: 11,
            additional_info: "Очный тур в Санкт-Петербурге",
            status: ApplicationStatus.Approved,
        },
        {
            id: "3",
            name: "Этап апелляции по информатике",
            start_date: "2025-11-15T10:00:00Z",
            end_date: "2025-11-15T13:00:00Z",
            previous_event_id: "2",
            subject: 3,
            class_number: 9,
            additional_info: "Подача заявлений до 14 ноября",
            status: ApplicationStatus.Rejected,
        },
    ];

    useEffect(() => {
        async function fetchAppeal(userId: string) {
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

        // fetchAppeal("")
        setLoading(false)
        setEvents(mockEvents)
    }, [])


    if (loading) return <div>Загрузка...</div>

    if (events.length === 0) return <div>Вы не подавали аппеляции</div>

    function footer(status: number, event: ApplicationEvent): ReactNode {
        return (
            <div className="d-flex flex-column justify-content-between h-100">
                <button className="btn btn-primary mb-2" onClick={() => navigate(`/profile/appeals/${1}/appeal-view`)}>Результаты</button>
                <div className="text-end">
                    <b>Статус апелляции  <StatusIcon status={status} /></b>
                </div>
            </div>
        )
    }

    return (
        <div>
            {events.map((event) => (
                <ApplicationEventCart key={event.id} event={event} footer={footer(event.status, event)} />
            ))}
        </div>
    );
};




export default AppealTab;