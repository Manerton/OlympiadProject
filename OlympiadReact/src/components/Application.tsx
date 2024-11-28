import React from "react";
import { Accordion,Badge,Button } from "react-bootstrap";
import { Application } from "../types/application.ts"

interface ApplicationCardProps {
  application: Application;
}

const ApplicationCard: React.FC<ApplicationCardProps> = ({ application }) => {
    const getStatusLabel = () => {
        if (application.status === true) {
          return (
            <Badge bg="success" className="px-2 py-1">
              Одобрено
            </Badge>
          );
        }
        if (application.status === false) {
          return (
            <Badge bg="danger" className="px-2 py-1">
              Отклонено
            </Badge>
          );
        }
        return (
          <Badge bg="primary" className="px-2 py-1">
            Не обработано
          </Badge>
        );
      };

      /* return (
        <Card className="p-3 mb-3 shadow-sm" style={{ maxWidth: "400px" }}>
          <Card.Body className="p-2">
            <Card.Title className="text-center mb-2">
              Заявка №{application.applicationID}
            </Card.Title>
            <hr className="my-2" />
            <div className="text-start">
              <Card.Text>
                <strong>Событие Имя:</strong> {application.eventID}
              </Card.Text>
              <Card.Text>
                <strong>Статус:</strong> {getStatusLabel()}
              </Card.Text>
              <Card.Text>
                <strong>Дата подачи:</strong>{" "}
                {new Date(application.submittedAt).toLocaleString()}
              </Card.Text>
              <Card.Text>
                <strong>Последнее обновление:</strong>{" "}
                {new Date(application.updatedAt).toLocaleString()}
              </Card.Text>
            </div>
          </Card.Body>
        </Card>
      ); */
      return (
        <Accordion className="mb-3 shadow-sm">
          <Accordion.Item eventKey="0">
          <Accordion.Header>
            Заявка №{application.applicationID}
            <span className="ms-2">{getStatusLabel()}</span>
          </Accordion.Header>
            <Accordion.Body>
              <div className="text-start">
                <p>
                  <strong>Событие:</strong> {application.eventName}
                </p>
                <p>
                  <strong>Место:</strong> {application.eventLocation}
                </p>
                <p>
                  <strong>Дата проведения:</strong> {" "}
                  {new Date(application.eventDate).toLocaleString()}
                </p>
                <p>
                  <strong>Дата подачи:</strong>{" "}
                  {new Date(application.submittedAt).toLocaleString()}
                </p>
                <p>
                  <strong>Последнее обновление:</strong>{" "}
                  {new Date(application.updatedAt).toLocaleString()}
                </p>
                <p className="d-flex justify-content-center mt-3">
                    <Button variant="info" size="sm">
                      Подробнее
                    </Button>
                </p>
                </div>
            </Accordion.Body>
          </Accordion.Item>
        </Accordion>
      );
    };

export default ApplicationCard;
