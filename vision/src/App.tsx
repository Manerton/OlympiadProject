import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import UserIndex from './components/olymp-admin/user/index.tsx';
import UserShow from './components/olymp-admin/user/show.tsx';

function App() {
    return (
        <Router>
            <div className="App">
                <Routes>
                    <Route path="/olymp-admin/user/index" element={<UserIndex />} />
                    <Route path="/olymp-admin/user/show/:id" element={<UserShow />} />
                </Routes>
            </div>
        </Router>
    );
}
export default App;