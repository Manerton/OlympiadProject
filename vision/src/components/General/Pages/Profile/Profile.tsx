import { useState } from "react";
import { UserRole } from "../../../../dictionary/role";
import ApplicationEventTab from "./Tabs/ApplicationEvents";
import HistoryTab from "./Tabs/History";
import AppealTab from "./Tabs/Appeal/Appeal";
import type { ApplicationEvent } from "../../../types/event";
import ResultByEvent from "./Tabs/ResultByEvent";
import ResultAppeal from "./Tabs/Appeal/ResultAppeal";

// общий тип для таба
type TabItem = {
  label: string;
  component?: React.ComponentType<any>;
  detailComponent?: React.ComponentType<{ event: ApplicationEvent; onBack: () => void }>;
};


const ProfileMainPage: React.FC = () => {

    const role = 2

    const tabsByRole: Record<number, TabItem[]> = {
        [UserRole.Participant]: [
            { label: "Достжения" },
            { label: "Олимпиады", component: ApplicationEventTab},
            { label: "История", component: HistoryTab, detailComponent: ResultByEvent},
            { label: "Апелляции", component: AppealTab, detailComponent: ResultAppeal },
        ],
        [UserRole.Judge]: [
            { label: "История" },
            { label: "Олимпиады" },
        ]
    }

    const tabs = tabsByRole[role] || []
    const [activeTab, setActiveTab] = useState(0)
    const [selectedEvent, setSelectedEvent] = useState<ApplicationEvent | null>(null);


    const ActiveComponent = tabs[activeTab].component;
    const DetailsComponent = tabs[activeTab].detailComponent

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
                {DetailsComponent && selectedEvent ? (
                    <DetailsComponent event={selectedEvent} onBack={() => setSelectedEvent(null)} />
                ) : ActiveComponent ? (
                    // пробрасываем селектор события только туда, где он нужен
                    <ActiveComponent onSelectEvent={setSelectedEvent} />
                ) : null}
            </div>
        </div>
    );
};

export default ProfileMainPage
