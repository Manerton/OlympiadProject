import React from 'react'
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom'
import Header from './components/header.tsx'
import ApplicationsPage from './components/ApplicationsList.tsx'
import AttendancePage from './components/AttendaceList.tsx'
import './App.css'
import ProfilePage from './components/Profile.tsx'
import Login from './components/Login.tsx'
import SubStagePage from './components/eventsServiceComponents/pages/subStagePage.tsx'
import OlympiadStagesPage from './components/eventsServiceComponents/pages/olympiadStagesPage.tsx'
import RegionalStagesPage from './components/eventsServiceComponents/pages/regionalStagesPage.tsx'
import OlympiadsPage from './components/eventsServiceComponents/pages/olympiadsPage.tsx'
import ApplicationOrganizerPage from './components/ApplicationOrganizerPage.tsx' 
import { RoleProvider } from "./components/RoleContext";
import OlympiadClassPage from './components/eventsServiceComponents/pages/olympiadClass.tsx'

function App() {


  return (
    <Router>
      <RoleProvider>
        <Header />
      
      <div className="container mt-4">
        <Routes>
          <Route path="/events" element={<RegionalStagesPage />} />
          <Route path="/olympiads/:id" element={<OlympiadsPage />} />
          <Route path="/olympiad-class/:id" element={<OlympiadClassPage />}/>
          <Route path="/olympiad-stages/:id" element={<OlympiadStagesPage />} />
          <Route path="/sub-stage/:id" element={<SubStagePage />}/>
          <Route path="/applications/user/:id" element={<ApplicationsPage />} />
          <Route path="/applications/event/:id" element={<ApplicationOrganizerPage />} />
          <Route path="/attendance" element={<AttendancePage />} />
          <Route path="/profile" element={<ProfilePage />} />
          <Route path="/login" element={<Login />} />
        </Routes>
      </div>
      </RoleProvider>

    </Router>
  )
}

export default App
