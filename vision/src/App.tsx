import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import UserIndex from './components/olymp-admin/user/index.tsx';
import UserShow from './components/olymp-admin/user/show.tsx';
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
                </Routes>
            </div>
            </Layout>
        </Router>
    );
}
export default App;