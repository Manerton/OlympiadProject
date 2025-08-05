import React from 'react';
import { Container, Row, Col } from 'react-bootstrap';

const Footer: React.FC = () => (
  <footer className="pt-5 border bg-body-tertiary">
    {/* Основной блок */}
    <Container fluid="lg">
      <Row className="gy-4">
        {/* Информация */}
        <Col md={4}>
          <h6 className="text-uppercase fw-bold mb-3 text-center">Информация</h6>
          <p className="small">
            Система для проведения Всероссийской олимпиады школьников. Она создана для упрощения процессов, помочь талантливым школьникам раскрыть свой потенциал. <br/>
            С помощью этой системы можно зарегистрироваться, подать заявку на участие, узнать результаты и получить всю необходимую информацию в одном месте.
          </p>
        </Col>

        {/* Представители */}
        <Col md={4} className="text-center">
          <h6 className="text-uppercase fw-bold mb-3">Представители</h6>
          <div className="d-flex justify-content-center align-items-center gap-3">
            <img src="src/assets/images/v51_9.png" alt="Всероссийская олимпиады" style={{ maxHeight: '120px' }} />
            
          </div>
          <img src="src/assets/images/v49_8.png" alt="РШТ" style={{ maxHeight: '70px' }} />
        </Col>

        {/* Контактная информация */}
        <Col md={4}>
          <h6 className="text-uppercase fw-bold mb-3 text-center">Контактная информация</h6>
          <address className="small mb-3">
            <strong>г. Москва</strong><br/>
            ул. Жуковского, д.16<br/>
            Телефон: <a href="tel:+74956255980" className="text-white">8 (495) 625-59-80</a><br/>
            <a href="mailto:Fcod@edu.gov.ru" className="text-white">Fcod@edu.gov.ru</a>
          </address>
          <address className="small">
            <strong>г. Астрахань</strong><br/>
            ул. Анри Барбюса, 7<br/>
            Телефон: <a href="tel:+78512442428" className="text-white">+7 8512 442-428</a><br/>
            <a href="mailto:schooltech@astrobI.ru" className="text-white">schooltech@astrobI.ru</a>
          </address>
        </Col>
      </Row>

      {/* Социальные иконки */}
      <Row className="mt-4">
        <Col className="text-center">
          <div className="d-flex justify-content-center align-items-center gap-4">
            <a href="#" aria-label="ВКонтакте">
              <img src= "src/assets/images/v51_12.png" alt="VK" style={{ height: '24px' }} />
            </a>
            <a href="#" aria-label="Школьная сеть">
              <img src="src/assets/images/v51_13.png" alt="School" style={{ height: '24px' }} />
            </a>
            <a href="#" aria-label="Telegram">
              <img src="src/assets/images/v51_11.png" alt="Telegram" style={{ height: '24px' }} />
            </a>
          </div>
        </Col>
      </Row>
    </Container>

    {/* Нижняя полоса */}
    <div className="bg-secondary text-center py-3 mt-5">
      <small>© 2025 Все права защищены</small>
    </div>
  </footer>
);

export default Footer;

