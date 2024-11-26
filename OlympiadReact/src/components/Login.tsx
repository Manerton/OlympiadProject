import React from "react";
import { Form, Button, Container, Row, Col, Alert } from "react-bootstrap";
import { Formik } from "formik";
import * as Yup from "yup";

interface LoginFormValues {
  username: string;
  password: string;
}

const Login: React.FC = () => {
  const initialValues: LoginFormValues = {
    username: "",
    password: "",
  };

  const validationSchema = Yup.object({
    username: Yup.string().email("Неверный формат email").required("Email обязателен"),
    password: Yup.string().required("Пароль обязателен").min(6, "Пароль должен быть не менее 6 символов"),
  });

  const handleSubmit = (values: LoginFormValues) => {
    // Здесь вы можете обрабатывать авторизацию
    console.log(values);
    // Пример: можно сделать запрос на сервер
  };

  return (
    <Container className="d-flex justify-content-center align-items-center" style={{ height: "50vh" }}>
     
        <Col xs={12} md={6} lg={4}>
          <div className="text-center mb-4">
            <h2>Авторизация</h2>
          </div>

          <Formik
            initialValues={initialValues}
            validationSchema={validationSchema}
            onSubmit={handleSubmit}
          >
            {({
              handleSubmit,
              handleChange,
              values,
              errors,
              touched,
            }) => (
              <Form onSubmit={handleSubmit}>
                <Form.Group className="mb-3" controlId="username">
                  <Form.Label>Электронная почта</Form.Label>
                  <Form.Control
                    type="email"
                    placeholder="Введите email"
                    name="username"
                    value={values.username}
                    onChange={handleChange}
                    isInvalid={touched.username && !!errors.username}
                  />
                  <Form.Control.Feedback type="invalid">
                    {errors.username}
                  </Form.Control.Feedback>
                </Form.Group>

                <Form.Group className="mb-3" controlId="password">
                  <Form.Label>Пароль</Form.Label>
                  <Form.Control
                    type="password"
                    placeholder="Введите пароль"
                    name="password"
                    value={values.password}
                    onChange={handleChange}
                    isInvalid={touched.password && !!errors.password}
                  />
                  <Form.Control.Feedback type="invalid">
                    {errors.password}
                  </Form.Control.Feedback>
                </Form.Group>

                {/* Статус авторизации */}
                <Alert variant="danger" className="d-none">
                  Неверный логин или пароль
                </Alert>

                <Button variant="primary" type="submit" className="d-block w-100">
                Войти
                </Button>
              </Form>
            )}
          </Formik>
        </Col>
     
    </Container>
  );
};

export default Login;
