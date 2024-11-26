import { BrowserRouter as Router, Route, Routes } from 'react-router-dom'
import Header from './components/header'
import ApplicationsPage from './components/ApplicationsList'
import AttendancePage from './components/AttendaceList'
import './App.css'
import ProfilePage from './components/Profile'
import Login from './components/Login'

function App() {
 

  return (
    <Router>
    <Header />
    <div className="container mt-4">
    <Routes>
        <Route path="/applications"  element={<ApplicationsPage />} />
        <Route path="/attendance"  element={<AttendancePage />} />
        <Route path="/profile"  element={<ProfilePage />} />
        <Route path="/login"  element={<Login />} />
    </Routes>
    </div>
  </Router>
  )
}

export default App
