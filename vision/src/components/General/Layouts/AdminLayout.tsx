import React from 'react';
import AdminSidebar from '../../Admin/AdminComponents/AdminSidebar.tsx';
import Footer from './Footer';
import { Container } from 'react-bootstrap';
import { Outlet } from 'react-router-dom';

const AdminLayout: React.FC = () => {
  return (
    <div className="d-flex min-vh-100">
      <AdminSidebar />

      <div className="flex-grow-1 d-flex flex-column">
        <Container fluid="lg" className="flex-grow-1">
          <main className="py-4">
            <Outlet />
          </main>
        </Container>
        <Footer />
      </div>
    </div>
  );
};

export default AdminLayout;
