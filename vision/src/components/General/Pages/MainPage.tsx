// MainPage.tsx
import React from 'react';
import { Container, Row, Col, Button, Form, Accordion } from 'react-bootstrap';

const MainPage: React.FC = () => {
  return (
    <>
      {/* Hero Section */}
      <Row className="g-0 border border-1 rounded justify-content-center align-items-center">
           <Col md={6}>
                <div
                  className="d-flex flex-column justify-content-center align-items-center h-100 text-center p-4"
                  style={{ backgroundColor: 'rgba(177, 172, 172, 0.5)', minHeight: '70dvh' }}
                >
                  <h1 className="display-4 fw-bold text-center">
                      Всероссийская олимпиада школьников <br />
                      <span
                        style={{
                          background: 'linear-gradient(to right, #1494D4, #70FF99)',
                          WebkitBackgroundClip: 'text',
                          WebkitTextFillColor: 'transparent',
                          display: 'inline-block'
                        }}
                      >
                        в твоём регионе!
                      </span>
                    </h1>

                  <p className="lead">
                    Упрощение процессов, помощь талантливым школьникам раскрыть свой потенциал.
                  </p>
                </div>
              </Col>

           <Col md={6}>
            <div
                  className="position-relative d-flex align-items-center justify-content-center"
                  style={{ minHeight: '70dvh', backgroundColor: 'rgba(177, 172, 172, 0.5)' }}
                >
                  <img
                    src="src/assets/images/v12_13.png"
                    alt="Background"
                    className="position-absolute w-100 h-100"
                    style={{ objectFit: 'cover', zIndex: 0 }}
                  />
            </div>

           </Col> 
      </Row>
      

      {/* How It Works Section */}
      <section className="mt-4 py-3 text-center border border-1">
        <Container fluid="lg" className="mx-auto">
          <h2 className="mb-4">Как проходит олимпиада?</h2>
          <Row className="g-4">
            <Col md={3}>
              <img src="src/assets/Point1.svg" className="img-fluid mb-3" alt="Присутствие участника" />
              <h5>Присутствие участника</h5>
              <p>Организатор отмечает явку участника в системе сопровождения ВсОШ. Участник подтверждает свою явку.</p>
            </Col>
            <Col md={3}>
              <img src="src/assets/Point2.svg" className="img-fluid mb-3" alt="Правила" />
              <h5>Правила</h5>
              <p>Перед началом тура участники проходят инструктаж, где им разъясняют правила проведения олимпиады, требования к оформлению работы, допустимые материалы и оборудование, правила поведения и последствия нарушения правил.</p>
            </Col>
            <Col md={3}>
              <img src="src/assets/Point3.svg" className="img-fluid mb-3" alt="Написание работы" />
              <h5>Написание работы</h5>
              <p>Участники самостоятельно выполняют задания в течение отведенного времени. После работы сдаются организаторам.</p>
            </Col>
            <Col md={3}>
              <img src="src/assets/Point4.svg" className="img-fluid mb-3" alt="Оглашение результатов" />
              <h5>Оглашение результатов</h5>
              <p>После проверки работ жюри результаты олимпиады публикуются в личных кабинетах участников и на сайте ВсОШ.</p>
            </Col>
          </Row>
        </Container>
      </section>

    {/* How It  */}
      <section className="mt-4 mb-4 py-5 text-center border border-1">
        <Container fluid="lg">
          <h2 className="mb-5">Как принять участие в олимпиаде?</h2>
          <Row
            className="gx-0 justify-content-center"
            style={{ overflow: 'visible' }}
          >

            {/* Шаг 1 */}
            <Col md={4}>
              <div className="p-4 rounded-4 h-100" style={{ backgroundColor: '#d8ffe4' }}>
                <div
                  className="d-flex flex-column align-items-center position-relative mb-4"
                  style={{ height: '6rem' }}
                >
                  <div
                    className="fw-bold text-success position-absolute top-50 start-50"
                    style={{
                      fontSize: '9rem',
                      opacity: 0.2,
                      transform: 'translate(-50%, -50%)',
                      zIndex: 0
                    }}
                  >
                    01
                  </div>
                  <h5 className="position-relative z-1 fw-bold text-dark" style={{ marginTop: '3rem' }}>
                    Регистрация
                  </h5>
                </div>
                <p className="text-success px-3" style={{ textAlign: 'center' }}>
                  Создайте личный кабинет на платформе. Введите необходимые данные и подтвердите свою учётную запись.
                </p>
              </div>
            </Col>

            {/* Шаг 2 — «выпирающая» карточка */}
            <Col
              md={4}
              style={{
                position: 'relative',
                zIndex: 1,
                marginTop: '-1rem',
                marginBottom: '-1rem'
              }}
            >
              <div className="p-4 rounded-4 h-100" style={{ backgroundColor: '#dfeeff' }}>
                <div
                  className="d-flex flex-column align-items-center position-relative mb-4"
                  style={{ height: '6rem' }}
                >
                  <div
                    className="fw-bold text-primary position-absolute top-50 start-50"
                    style={{
                      fontSize: '9rem',
                      opacity: 0.2,
                      transform: 'translate(-50%, -50%)',
                      zIndex: 0
                    }}
                  >
                    02
                  </div>
                  <h5 className="position-relative z-1 fw-bold text-dark" style={{ marginTop: '3rem' }}>
                    Подача заявки
                  </h5>
                </div>
                <p className="text-primary px-3" style={{ textAlign: 'center' }}>
                  Выберите интересующую вас олимпиаду из списка доступных и подайте заявку на участие. Дождитесь одобрения заявки организатором.
                </p>
              </div>
            </Col>

            {/* Шаг 3 */}
            <Col md={4}>
              <div className="p-4 rounded-4 h-100" style={{ backgroundColor: '#ffe4e9' }}>
                <div
                  className="d-flex flex-column align-items-center position-relative mb-4"
                  style={{ height: '6rem' }}
                >
                  <div
                    className="fw-bold text-danger position-absolute top-50 start-50"
                    style={{
                      fontSize: '9rem',
                      opacity: 0.2,
                      transform: 'translate(-50%, -50%)',
                      zIndex: 0
                    }}
                  >
                    03
                  </div>
                  <h5 className="position-relative z-1 fw-bold text-dark" style={{ marginTop: '3rem' }}>
                    Результаты
                  </h5>
                </div>
                <p className="text-danger px-3" style={{ textAlign: 'center' }}>
                  Примите участие в выбранной олимпиаде. После завершения этапов вы сможете просмотреть свои результаты, узнать о возможности подачи апелляции.
                </p>
              </div>
            </Col>

          </Row>
        </Container>
      </section>




      {/* FAQ Section */}
      <section className="mt-4 mb-4  py-3 text-center border rounded">
        <Container fluid="lg" className="mx-auto">
          <h2 className="mb-4">Часто задаваемые вопросы</h2>
          <Accordion defaultActiveKey="0">
            <Accordion.Item eventKey="0">
              <Accordion.Header>Как подать апелляцию?</Accordion.Header>
              <Accordion.Body>После публикации результатов нажмите кнопку «Подать апелляцию» в личном кабинете.</Accordion.Body>
            </Accordion.Item>
            <Accordion.Item eventKey="1">
              <Accordion.Header>Что делать, если не согласен с результатами?</Accordion.Header>
              <Accordion.Body>Вы можете подать апелляцию в течение 5 дней после объявления результатов.</Accordion.Body>
            </Accordion.Item>
            <Accordion.Item eventKey="2">
              <Accordion.Header>Когда и где будут проходить этапы олимпиады?</Accordion.Header>
              <Accordion.Body>Информация о датах и местах проведения этапов олимпиады будет опубликована на сайте и в личных кабинетах участников.</Accordion.Body>
            </Accordion.Item>
            <Accordion.Item eventKey="3">
              <Accordion.Header>Можно ли участвовать в олимпиаде по нескольким предметам?</Accordion.Header>
              <Accordion.Body>Да, участники могут подавать заявки на участие в олимпиадах по нескольким предметам.</Accordion.Body>
            </Accordion.Item>
            <Accordion.Item eventKey="4">
              <Accordion.Header>Что делать, если я забыл пароль?</Accordion.Header>
              <Accordion.Body>На странице входа нажмите «Забыли пароль?» и следуйте инструкциям для восстановления пароля.</Accordion.Body>
            </Accordion.Item>
            <Accordion.Item eventKey="5">
              <Accordion.Header>Как подать заявку на участие в олимпиаде?</Accordion.Header>
              <Accordion.Body>После регистрации и входа в личный кабинет выберите интересующую вас олимпиаду из списка доступных и нажмите кнопку “Подать заявку”.</Accordion.Body>
            </Accordion.Item>
          </Accordion>
        </Container>
      </section>

      {/* Contact Section */}
      <section className="mt-4 mb-4 py-3 border rounded">
        <Container fluid="lg" className="mx-auto">
          <h2 className="mb-4 text-center">Свяжитесь с нами</h2>
          <Form>
            <Row className="g-3">
              <Col md={6}>
                <Form.Group controlId="contactName">
                  <Form.Label>Имя</Form.Label>
                  <Form.Control type="text" placeholder="Ваше имя" />
                </Form.Group>
              </Col>
              <Col md={6}>
                <Form.Group controlId="contactEmail">
                  <Form.Label>Почта</Form.Label>
                  <Form.Control type="email" placeholder="Ваш email" />
                </Form.Group>
              </Col>
              <Col xs={12}>
                <Form.Group controlId="contactMessage">
                  <Form.Label>Сообщение</Form.Label>
                  <Form.Control as="textarea" rows={4} />
                </Form.Group>
              </Col>
              <Col xs={12} className="text-center">
                <Button type="submit">Отправить</Button>
              </Col>
            </Row>
          </Form>
        </Container>
      </section>
    </>
  );
};

export default MainPage;