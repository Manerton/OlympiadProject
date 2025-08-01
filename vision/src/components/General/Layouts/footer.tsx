import React from 'react';
import { Container } from 'react-bootstrap';

const Footer: React.FC = () => (
  <footer className="py-3 bg-dark text-white text-center">
    <Container>
      <p className="mb-0">© 2025 Все права защищены</p>
      <small>г. Москва, ул. Жуковского, д.16 | Тел: 8 (495) 625-59-80</small>
    </Container>
  </footer>
);

export default Footer;
