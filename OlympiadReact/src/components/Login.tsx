import React from "react";
import { Form, Button, Container, Row, Col, Alert } from "react-bootstrap";
import { Formik } from "formik";
import * as Yup from "yup";
import { redirect } from "react-router-dom";

interface LoginFormValues {
  email: string;
  password: string;
}

const Login: React.FC = () => {
  const initialValues: LoginFormValues = {
    email: "",
    password: "",
  };

  const validationSchema = Yup.object({
    email: Yup.string().email("Неверный формат email").required("Email обязателен"),
    password: Yup.string().required("Пароль обязателен").min(4, "Пароль должен быть не менее 4 символов"),
  });

  const handleSubmit = async (values: LoginFormValues) => {
    // Включите обработку логики авторизации
    console.log(values); // values содержит данные формы (например, email, password)
  
    const { email, password } = values;
  
    try {
      // Делаем POST запрос с данными формы
      const response = await fetch('http://localhost:8081/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json', // Указываем, что отправляем JSON
        },
        body: JSON.stringify({ email, password }), // Отправляем данные формы в теле запроса
        credentials: 'include', // Включает cookies
      });
  
      // Проверяем статус ответа
      if (!response.ok) {
        // Если ответ не успешный (например, статус 400 или 500)
        throw new Error('Авторизация не удалась');
      }
      
      
      // Ответ успешный, можем обработать данные
      console.log('Авторизация успешна:');
  
      
    } catch (error) {
      // Обрабатываем ошибки
      console.error('Ошибка авторизации:', error);
      alert('Ошибка при авторизации. Попробуйте снова.');
    }
    redirect("/profile")
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
                <Form.Group className="mb-3" controlId="email">
                  <Form.Label>Электронная почта</Form.Label>
                  <Form.Control
                    type="email"
                    placeholder="Введите email"
                    name="email"
                    value={values.email}
                    onChange={handleChange}
                    isInvalid={touched.email && !!errors.email}
                  />
                  <Form.Control.Feedback type="invalid">
                    {errors.email}
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
