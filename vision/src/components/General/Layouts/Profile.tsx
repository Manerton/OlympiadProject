import React from 'react';
import AdminSidebar from '../../Admin/AdminComponents/AdminSidebar.tsx';
import { Container } from 'react-bootstrap';
import { NavLink, Outlet } from 'react-router-dom';
import Header from './header.tsx';
import Footer from './Footer.tsx';
import { UserRole } from '../../../dictionary/role.tsx';
import { useAuth } from '../../Helpers/AuthContext.tsx';

const ProfileLayout: React.FC = () => {

    const tabsByRole: Record<number, { label: string; path: string }[]> = {
        [UserRole.Participant]: [
            { label: "Достижения", path: "achievements" },
            { label: "Олимпиады", path: "applications" },
            { label: "История", path: "history" },
            { label: "Апелляции", path: "appeals" },
        ],
        [UserRole.Judge]: [
            { label: "История", path: "history" },
            { label: "Олимпиады", path: "applications" },
        ],
    };

    const { user } = useAuth();

    const role = UserRole.Participant; // например, из контекста или стора
    const tabs = tabsByRole[role];

    return (
        <div>
            <Header />

            <div className="d-flex flex-column min-vh-100">
                <div className="container">
                    {/* Информация о пользователе */}
                    <div className="d-flex align-items-center m-4">
                        <img src="" alt="" className="roumded-circle me-3" width={80} height={80} />
                        <div>
                            <h4>Иванов Иван Иванович</h4>
                            <p className="text-muted">{user?.Email}</p>
                        </div>
                        <button className="btn btn-primary ms-auto">Редактировать</button>
                    </div>

                    <ul className="nav nav-tabs w-100">
                        {tabs.map((tab) => (
                            <li key={tab.path} className="nav-item flex-fill text-center">
                                <NavLink
                                    to={tab.path}
                                    className={({ isActive }) =>
                                        "nav-link" + (isActive ? " active" : "")
                                    }
                                >
                                    {tab.label}
                                </NavLink>
                            </li>
                        ))}
                    </ul>

                    <div className="tab-content mt-3">
                        <Outlet />
                    </div>
                </div>
            </div>
        </div>
    );
};

export default ProfileLayout;
