import React from 'react';
import { NavLink, Outlet, useNavigate } from 'react-router-dom';
import Header from './header';
import Footer from './footer';
import { UserRole } from '../../../dictionary/role';
import { useAuth } from '../../Helpers/AuthContext';
import UserInfoBlock from '../../Common/UserInfo';

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

    const navigate = useNavigate()

    const actions = (
        <>
            <button className="btn btn-primary" onClick={() => navigate("/profile/edit")}>Подробнее</button>
        </>
    );

    return (
        <div>
            <Header />

            <div className="d-flex flex-column min-vh-100">
                <div className="container">
                    {/* Информация о пользователе */}
                    <UserInfoBlock email={user?.Email} actions={actions}/>

                    {/* Навигационные вкладки */}
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
