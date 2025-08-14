import { useState } from "react";
import { UserRole } from "../../../../dictionary/role";


const ProfileMainPage: React.FC = () => {

    const role = 2

    const tabsByRole = {
        [UserRole.Participant]: [ 
            {"label": "Достжения", "component": ""},
            {"label": "Олимпиады", "component": ""},
            {"label": "История", "component": ""},
            {"label": "Апелляцит", "component": ""},
        ],
        [UserRole.Judge]: [
            {"label": "История", "component": ""},
            {"label": "Олимпиады", "component": ""},
        ]
    }

    const tabs = tabsByRole[role] || []
    const [acriveTab, setActiveTab] = useState(0)

    return (
        <div className="container mt-4">
            {/* Информация о пользователе */}
            <div className="d-flex align-items-center mb-4">
                <img src="" alt="" className="roumded-circle me-3" width={80} height={80}/>
                <div>
                    <h4>FIO</h4>
                    <p className="text-muted">EMAIL</p>
                </div>
                <button className="btn btn-primary ms-auto">Редактировать</button>
            </div>
            {/* Табы */}
            <ul className="nav nav-tabs">
                {tabs.map((tab, i) => (
                    <li key={i} className="nav-item">
                        <button  
                        className={`nav-link ${i === acriveTab ? "active" : ""}`}
                            onClick={() => setActiveTab(i)}>
                            {tab.label}
                        </button>
                    </li>
                ))}
            </ul>

            {/* Содержимое */}
            <div className="tab-content p-3 border border-top-0">
                {tabs[acriveTab].component}
            </div>

        </div>
    );
};

export default ProfileMainPage