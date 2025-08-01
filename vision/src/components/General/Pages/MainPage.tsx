import React from 'react';
import { Container, Row, Col, Button, Form, Navbar, Nav, Accordion } from 'react-bootstrap';

const MainPage: React.FC = () => {
  return (
    <>
      <Container fluid className="p-0">
        {/* Главная секция (Hero) */}
        <section
          className="vh-100 d-flex align-items-center justify-content-center bg-light"
          style={{
            backgroundImage: "url('src/assets/images/v12_13.png')",
            backgroundSize: 'cover',
            backgroundPosition: 'center',
          }}
        >
          <div className="text-center text-white p-4" style={{ backgroundColor: 'rgba(0,0,0,0.5)', borderRadius: '1rem' }}>
            <h1 className="display-4 fw-bold">
              Всероссийская олимпиада школьников в твоём регионе!
            </h1>
            <p className="lead">
              Упрощение процессов, помощь талантливым школьникам раскрыть свой потенциал.
            </p>
          </div>
        </section>

            {/* Секция "Как проходит олимпиада?" */}
        <section className="py-5">
          <Container>
            <h2 className="mb-4">Как проходит олимпиада?</h2>
            <Row className="g-4">
              <Col md={3}>
               <img src="src/assets/Point1.svg" className="img-fluid mb-3" />
                <h5>Присутствие участника</h5>
                <p>Организатор отмечает явку участника в системе сопровождения ВсОШ. Участник подтверждает свою явку.</p>
              </Col>
              
              <Col md={3}>
               <img src="src/assets/Point2.svg" className="img-fluid mb-3" />
                <h5>Правила</h5>
                <p>Перед началом тура участники проходят инструктаж, где им разъясняют правила проведения олимпиады, требования к оформлению работы, допустимые материалы и оборудование, правила поведения и последствия нарушения правил.</p>
              </Col>
              
              <Col md={3}>
               <img src="src/assets/Point3.svg" className="img-fluid mb-3" />
                <h5>Написание работы</h5>
                <p>Участники самостоятельно выполняют задания в течение отведенного времени. После работы сдаются организаторам.</p>
              </Col>
              <Col md={3}>
              <img src="src/assets/Point4.svg" className="img-fluid mb-3" />
                <h5>Оглашение результатов</h5>
                <p>После проверки работ жюри результаты олимпиады публикуются в личных кабинетах участников и на сайте ВсОШ.</p>
              </Col>
            </Row>
          </Container>
        </section>


        {/* Секция "Как принять участие в олимпиаде?" */}
        <section className="py-5 bg-dark text-white">
          <Container>
            <h2 className="mb-4 text-center">Как принять участие в олимпиаде?</h2>
            <Row className="text-center g-4">
              <Col md={3}>
                <img src="src/assets/images/v51_11.png" alt="Регистрация" className="img-fluid mb-3" />
                <h5 className="fw-bold">Регистрация</h5>
                <p>Создайте личный кабинет на платформе. Введите необходимые данные и подтвердите свою учётную запись.</p>
              </Col>
              <Col md={3}>
                <img src="src/assets/images/v51_12.png" alt="Подача заявки" className="img-fluid mb-3" />
                <h5 className="fw-bold">Подача заявки</h5>
                <p>Выберите интересующую вас олимпиаду из списка доступных и подайте заявку на участие. Дождитесь одобрения заявки организатором.</p>
              </Col>
              <Col md={3}>
                <img src="src/assets/images/v51_13.png" alt="Написание работы" className="img-fluid mb-3" />
                <h5 className="fw-bold">Написание работы</h5>
                <p>Выполните задания в личном кабинете в отведенное время.</p>
              </Col>
              <Col md={3}>
                <img src="src/assets/images/v51_14.png" alt="Результаты" className="img-fluid mb-3" />
                <h5 className="fw-bold">Результаты</h5>
                <p>Примите участие в выбранной олимпиаде. После завершения этапов вы сможете просмотреть свои результаты, узнать о возможности подачи апелляции.</p>
              </Col>
            </Row>
          </Container>
        </section>

      

        {/* Секция "Часто задаваемые вопросы" */}
        <section className="py-5">
          <Container>
            <h2 className="mb-4">Часто задаваемые вопросы</h2>
            <Accordion defaultActiveKey="0">
              <Accordion.Item eventKey="0">
                <Accordion.Header>Как подать апелляцию?</Accordion.Header>
                <Accordion.Body>
                  После публикации результатов нажмите кнопку «Подать апелляцию» в личном кабинете.
                </Accordion.Body>
              </Accordion.Item>
              <Accordion.Item eventKey="1">
                <Accordion.Header>Что делать, если не согласен с результатами?</Accordion.Header>
                <Accordion.Body>
                  Вы можете подать апелляцию в течение 5 дней после объявления результатов.
                </Accordion.Body>
              </Accordion.Item>
              <Accordion.Item eventKey="2">
                <Accordion.Header>Когда и где будут проходить этапы олимпиады?</Accordion.Header>
                <Accordion.Body>
                  Информация о датах и местах проведения этапов олимпиады будет опубликована на сайте и в личных кабинетах участников.
                </Accordion.Body>
              </Accordion.Item>
              <Accordion.Item eventKey="3">
                <Accordion.Header>Можно ли участвовать в олимпиаде по нескольким предметам?</Accordion.Header>
                <Accordion.Body>
                  Да, участники могут подавать заявки на участие в олимпиадах по нескольким предметам.
                </Accordion.Body>
              </Accordion.Item>
              <Accordion.Item eventKey="4">
                <Accordion.Header>Что делать, если я забыл пароль?</Accordion.Header>
                <Accordion.Body>
                  На странице входа нажмите «Забыли пароль?» и следуйте инструкциям для восстановления пароля.
                </Accordion.Body>
              </Accordion.Item>
              <Accordion.Item eventKey="5">
                <Accordion.Header>Как подать заявку на участие в олимпиаде?</Accordion.Header>
                <Accordion.Body>
                  После регистрации и входа в личный кабинет выберите интересующую вас олимпиаду из списка доступных и нажмите кнопку “Подать заявку”.
                </Accordion.Body>
              </Accordion.Item>
            </Accordion>
          </Container>
        </section>

        {/* Секция "Свяжитесь с нами" */}
        <section className="py-5 bg-light">
          <Container>
            <h2 className="mb-4">Свяжитесь с нами</h2>
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
                <Col xs={12}>
                  <Button type="submit">Отправить</Button>
                </Col>
              </Row>
            </Form>
          </Container>
        </section>
      </Container>

      
    </>
  );
};

export default MainPage;