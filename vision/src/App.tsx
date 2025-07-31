import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import UserIndex from './components/olymp-admin/user/index.jsx';

function App() {
    return (
        <Router>
            <div className="App">
                <Routes>
                    {/* Этот маршрут будет работать для любого URL */}
                    <Route path="/olymp-admin/user/index" element={<UserIndex />} />
                </Routes>
            </div>
        </Router>
    );
}
export default App;