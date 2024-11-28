import { BrowserRouter as Router, Route, Routes } from 'react-router-dom'
import Header from './components/header'
import ApplicationsPage from './components/ApplicationsList'
import AttendancePage from './components/AttendaceList'
import './App.css'
import ProfilePage from './components/Profile'
import Login from './components/Login'
import SubStagePage from './components/eventsServiceComponents/pages/subStagePage'
import OlympiadStagesPage from './components/eventsServiceComponents/pages/olympiadStagesPage'
import RegionalStagesPage from './components/eventsServiceComponents/pages/regionalStagesPage'
import OlympiadsPage from './components/eventsServiceComponents/pages/olympiadsPage'

function App() {


  return (
    <Router>
      <Header />
      <div className="container mt-4">
        <Routes>
          <Route path="/events" element={<RegionalStagesPage />} />
          <Route path="/olympiads/:id" element={<OlympiadsPage />} />
          <Route path="/olympiad-stages/:id" element={<OlympiadStagesPage />} />
          <Route path="/sub-stage/:id" element={<SubStagePage />}/>
          <Route path="/applications" element={<ApplicationsPage />} />

          <Route path="/attendance" element={<AttendancePage />} />
          <Route path="/profile" element={<ProfilePage />} />
          <Route path="/login" element={<Login />} />
        </Routes>
      </div>
    </Router>
  )
}

export default App
