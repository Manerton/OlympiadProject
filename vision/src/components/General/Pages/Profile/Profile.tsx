import { useState } from "react";
import { UserRole } from "../../../../dictionary/role";
import ApplicationEventTab from "./Tabs/ApplicationEvents";
import HistoryTab from "./Tabs/History";
import AppealTab from "./Tabs/Appeal";
import type { ApplicationEvent } from "../../../types/event";

// общий тип для таба
type TabItem = {
    label: string;
    component?: React.ComponentType<any>; // можно передавать пропы
};


const ProfileMainPage: React.FC = () => {

    const role = 2

    const tabsByRole: Record<number, TabItem[]> = {
        [UserRole.Participant]: [
            { "label": "Достжения" },
            { "label": "Олимпиады", "component": ApplicationEventTab },
            { "label": "История", "component": HistoryTab },
            { "label": "Апелляции", "component": AppealTab },
        ],
        [UserRole.Judge]: [
            { "label": "История" },
            { "label": "Олимпиады" },
        ]
    }

    const tabs = tabsByRole[role] || []
    const [activeTab, setActiveTab] = useState(0)
    const [selectedEvent, setSelectedEvent] = useState<ApplicationEvent | null>(null);


    const ActiveComponent = tabs[activeTab].component;

    return (
        <div className="container mt-4">
            {/* Информация о пользователе */}
            <div className="d-flex align-items-center mb-4">
                <img src="" alt="" className="roumded-circle me-3" width={80} height={80} />
                <div>
                    <h4>Иванов Иван Иванович</h4>
                    <p className="text-muted">ivanovivan@main.com</p>
                </div>
                <button className="btn btn-primary ms-auto">Редактировать</button>
            </div>

            {/* Табы */}
            {!selectedEvent && (
                <ul className="nav nav-tabs nav-fill w-100">
                    {tabs.map((tab, i) => (
                        <li key={i} className="nav-item">
                            <button
                                className={`nav-link ${i === activeTab ? "active" : ""}`}
                                onClick={() => setActiveTab(i)}
                            >
                                {tab.label}
                            </button>
                        </li>
                    ))}
                </ul>
            )}

            {/* Содержимое */}
            <div className="tab-content p-3">
                {selectedEvent ? (
                    <EventDetail event={selectedEvent} onBack={() => setSelectedEvent(null)} />
                ) : ActiveComponent ? (
                    // пробрасываем селектор события только туда, где он нужен
                    <ActiveComponent onSelectEvent={setSelectedEvent} />
                ) : null}
            </div>
        </div>
    );
};

export default ProfileMainPage

function EventDetail({
    event,
    onBack,
}: {
    event: ApplicationEvent;
    onBack: () => void;
}) {
    return (
        <div>
            <button className="btn btn-link mb-3" onClick={onBack}>
                ← Назад
            </button>
            <h3 className="mb-3">{event.name}</h3>
            <p><strong>Дата начала:</strong> {new Date(event.start_date).toLocaleDateString()}</p>
            <p><strong>Дата окончания:</strong> {new Date(event.end_date).toLocaleDateString()}</p>
            <p>{event.additional_info}</p>
        </div>
    );
}
