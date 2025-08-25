import { use, useState } from "react";
import { UserRole } from "../../../../dictionary/role";
import ApplicationEventTab from "./Tabs/ApplicationEvents";
import HistoryTab from "./Tabs/History";
import AppealTab from "./Tabs/Appeal/Appeal";
import type { ApplicationEvent } from "../../../types/event";
import ResultByEvent from "./Tabs/ResultByEvent";
import ResultAppeal from "./Tabs/Appeal/ResultAppeal";
import { useAuth } from "../../../Helpers/AuthContext";
import AppealCreate from "./Tabs/Appeal/AppealCreate";
import type { Appeal } from "../../../types/appeal";

// общий тип для таба
type TabItem = {
  label: string;
  component?: React.ComponentType<any>;
  detailComponent?: React.ComponentType<{ event: ApplicationEvent; onBack: () => void }>;
  appealComponent?: React.ComponentType<{ onBack: () => void }>; //
};


const ProfileMainPage: React.FC = () => {

    const {user} = useAuth();

    const role = 2

    const tabsByRole: Record<number, TabItem[]> = {
        [UserRole.Participant]: [
            { label: "Достжения" },
            { label: "Олимпиады", component: ApplicationEventTab},
            { label: "История", component: HistoryTab, detailComponent: ResultByEvent, appealComponent: AppealCreate},
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
    const [appeal, setAppeal] = useState<Appeal | null>(null);

    const ActiveComponent = tabs[activeTab].component;
    const DetailsComponent = tabs[activeTab].detailComponent

    return (
        <div className="container mt-4">
            {/* Информация о пользователе */}
            <div className="d-flex align-items-center mb-4">
                <img src="" alt="" className="roumded-circle me-3" width={80} height={80} />
                <div>
                    <h4>Иванов Иван Иванович</h4>
                    <p className="text-muted">{user?.Email}</p>
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
                {appeal ? (
                    <AppealCreate onBack={() => setAppeal(null)} />
                ) : DetailsComponent && selectedEvent ? (
                    <DetailsComponent event={selectedEvent} onBack={() => setSelectedEvent(null)} />
                ) : ActiveComponent ? (
                    <ActiveComponent onSelectEvent={setSelectedEvent} onSelectAppeal={setAppeal} />
                ) : null}
            </div>

        </div>
    );
};

export default ProfileMainPage
