import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
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
import Header from './components/General/Layouts/header.tsx'
import ProfilePage from './components/General/Pages/Profile.tsx'
import MainPage from './components/General/Pages/MainPage.tsx'
import Layout from './components/General/Layouts/Layout.tsx'


function App() {
    return (
        <Router>
            <Layout>
            <div className="App">
                <Routes>
                     <Route path="/" element={<MainPage />} />
                    
                    <Route path="/olymp-admin/user/index" element={<UserIndex />} />
                    <Route path="/olymp-admin/user/show/:id" element={<UserShow />} />
                    <Route path="/olymp-admin/user/create" element={<UserCreate />} />
                    <Route path="/olymp-admin/user/edit/:id" element={<UserEdit />} />
                    <Route path="/profile" element={<ProfilePage />} />

                    <Route path="/olymp-admin/participant/index" element={<ParticipantIndex />} />
                    <Route path="/olymp-admin/participant/show/:id" element={<ParticipantShow />} />
                    <Route path="/olymp-admin/participant/create" element={<ParticipantCreate />} />
                    <Route path="/olymp-admin/participant/edit/:id" element={<ParticipantEdit />} />

                    <Route path="/olymp-admin/school/index" element={<SchoolIndex />} />
                    <Route path="/olymp-admin/school/show/:id" element={<SchoolShow />} />
                    <Route path="/olymp-admin/school/create" element={<SchoolCreate />} />
                    <Route path="/olymp-admin/school/edit/:id" element={<SchoolEdit />} />

                    <Route path="/olymp-admin/report/index" element={<ReportIndex />} />
                </Routes>
            </div>
            </Layout>
        </Router>
    );
}
export default App;