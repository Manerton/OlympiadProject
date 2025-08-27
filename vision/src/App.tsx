import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import UserIndex from './components/olymp-admin/user/index.tsx';
import UserShow from './components/olymp-admin/user/show.tsx';
import UserCreate from './components/olymp-admin/user/create.tsx';
import UserEdit from './components/olymp-admin/user/edit.tsx';
import ParticipantIndex from './components/olymp-admin/participant/index.tsx';
import ParticipantShow from './components/olymp-admin/participant/show.tsx';
import ParticipantCreate from './components/olymp-admin/participant/create.tsx';
import ParticipantEdit from './components/olymp-admin/participant/edit.tsx';
import SchoolIndex from './components/olymp-admin/school/index.tsx';
import SchoolCreate from './components/olymp-admin/school/create.tsx';
import SchoolShow from './components/olymp-admin/school/show.tsx';
import SchoolEdit from './components/olymp-admin/school/edit.tsx';
import ReportIndex from './components/olymp-admin/report/index.tsx';
import ApplicationIndex from './components/olymp-admin/application/index.tsx';
import ApplicationCreate from './components/olymp-admin/application/create.tsx';
import ApplicationShow from './components/olymp-admin/application/show.tsx';
import ApplicationEdit from './components/olymp-admin/application/edit.tsx';
import EventIndex from './components/olymp-admin/event/index.tsx';
import EventShow from './components/olymp-admin/event/show.tsx';
import EventAttendance from './components/olymp-admin/event/attendance.tsx';
import EventPrizeScore from './components/olymp-admin/event/prize-score.tsx'
import EventPoint from './components/olymp-admin/event/point.tsx'
import EventTask from './components/olymp-admin/event/task.tsx'
import Header from './components/General/Layouts/header.tsx'
import MainPage from './components/General/Pages/MainPage.tsx'
import AttendancePage from './components/General/Pages/AttendaceList.tsx'
import AdminPanel from './components/Admin/Pages/AdminPanel.tsx'
import Layout from './components/General/Layouts/Layout'
import AdminLayout from './components/General/Layouts/AdminLayout'
import AuthPage from './components/General/Pages/AuthForm.tsx'
import RegionalStagesPage from './components/General/Pages/Events/RegionalStages'
import OlympiadsPage from './components/General/Pages/Events/OlympiadsPage'
import RequireAuth from './components/Helpers/RequireAuth.tsx';
import ProfileLayout from './components/General/Layouts/Profile.tsx';
import HistoryTab from './components/General/Pages/Profile/Tabs/History.tsx';
import AppealTab from './components/General/Pages/Profile/Tabs/Appeal/Appeal.tsx';
import ApplicationEventTab from './components/General/Pages/Profile/Tabs/ApplicationEvents.tsx';
import ResultByEvent from './components/General/Pages/Profile/Tabs/ResultByEvent.tsx';
import AppealCreate from './components/General/Pages/Profile/Tabs/Appeal/AppealCreate.tsx';
import AppealView from './components/General/Pages/Profile/Tabs/Appeal/AppealView.tsx';
import EditProfile from './components/General/Pages/Profile/EditProfile.tsx';
import AppealList from './components/General/Pages/Profile/Tabs/Appeal/AppealList.tsx';


function App() {
    return (
        
        
        <Router>
            <div className="App">
                <Routes>
                    
                    <Route element={<Layout />}>
                        <Route path="/" element={<MainPage />} />
                        <Route path="/auth" element={<AuthPage />} />
                        {/* <Route
                            path="/profile"
                            element={
                                <RequireAuth>
                                <ProfileMainPage />
                                </RequireAuth>
                            }
                            /> */}
                        <Route path="/RegionalStages" element={<RegionalStagesPage />} />
                        <Route path= "OlympiadsPage/:id" element={<OlympiadsPage />} />
                    </Route>

                    <Route
                    element={
                        <RequireAuth allowedRoles={[1]}>
                        <AdminLayout />
                        </RequireAuth>
                    }
                    >
                    <Route path="/AdminPanel" element={<AdminPanel />} />
                    <Route path="/olymp-admin/user/index" element={<UserIndex />} />
                    <Route path="/olymp-admin/user/show/:id" element={<UserShow />} />
                    <Route path="/olymp-admin/user/create" element={<UserCreate />} />
                    <Route path="/olymp-admin/user/edit/:id" element={<UserEdit />} />
                    <Route path="/attendance" element={<AttendancePage />} />

                    <Route path="/olymp-admin/participant/index" element={<ParticipantIndex />} />
                    <Route path="/olymp-admin/participant/show/:id" element={<ParticipantShow />} />
                    <Route path="/olymp-admin/participant/create" element={<ParticipantCreate />} />
                    <Route path="/olymp-admin/participant/edit/:id" element={<ParticipantEdit />} />

                    <Route path="/olymp-admin/school/index" element={<SchoolIndex />} />
                    <Route path="/olymp-admin/school/show/:id" element={<SchoolShow />} />
                    <Route path="/olymp-admin/school/create" element={<SchoolCreate />} />
                    <Route path="/olymp-admin/school/edit/:id" element={<SchoolEdit />} />
                    <Route path="/olymp-admin/report/index" element={<ReportIndex />} />

                    <Route path="/olymp-admin/application/index" element={<ApplicationIndex />} />
                    <Route path="/olymp-admin/application/create" element={<ApplicationCreate />} />
                    <Route path="/olymp-admin/application/show/:id" element={<ApplicationShow />} />
                    <Route path="/olymp-admin/application/edit/:id" element={<ApplicationEdit />} />

                    <Route path="/olymp-admin/event/index" element={<EventIndex />} />
                    <Route path="/olymp-admin/event/show/:id" element={<EventShow />} />
                    <Route path="/olymp-admin/event/attendance/:id" element={<EventAttendance />} />
                    <Route path="/olymp-admin/event/prize-score/:id" element={<EventPrizeScore />} />
                    <Route path="/olymp-admin/event/task/:id" element={<EventTask />} />
                    <Route path="/olymp-admin/event/point/:id" element={<EventPoint />} />
                    </Route>

                    <Route path="/profile" element={<ProfileLayout/>}>
                        <Route index element={<Navigate to="history"/>}/>

                        <Route path="history" element={<HistoryTab/>}/>
                        <Route path="history/:eventId/result" element={<ResultByEvent/>}/>
                        <Route path="history/:eventId/appeal-create" element={<AppealCreate/>}/>

                        <Route path="applications" element={<ApplicationEventTab/>}/>

                        <Route path="appeals" element={<AppealTab/>}/>
                        <Route path="appeals/:appealId/appeal-view" element={<AppealView/>}/>
                        <Route path="appeals/:appealId/list" element={<AppealList/>}/>
                    </Route>

                    <Route  element={<Layout/>}>
                        <Route path="/profile/edit" element={<EditProfile/>}/>
                    </Route>

                </Routes>
            </div>
        </Router>
    );
}
export default App;