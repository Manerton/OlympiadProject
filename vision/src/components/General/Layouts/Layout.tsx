// Layout.tsx
import React, { ReactNode } from 'react';
import { Container } from 'react-bootstrap';
import Header from './Header';
import Footer from './Footer';

interface LayoutProps {
  children: ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  return (
    <>
      <Header />
      <div className="d-flex flex-column min-vh-100">
        <Container fluid="lg" className="flex-grow-1">
          <main className="py-4">
            {children}
          </main>
        </Container>
        <Footer />
      </div>
    </>
  );
};

export default Layout;




//  <div className="d-flex flex-column justify-content-center min-vh-100 min-vw-100">
//       <Header />
//        <Container fluid="lg" as="main" className="flex-grow-1 py-1">
//         {children}
//       </Container>
//       <Footer />
//     </div>