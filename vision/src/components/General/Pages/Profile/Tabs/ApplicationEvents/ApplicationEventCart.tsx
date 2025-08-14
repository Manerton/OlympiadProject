import { Card, Col } from "react-bootstrap";
import type { ApplicationEvent } from "../../../../../types/event"
import { getStatusLabel } from "../../../../../../dictionary/applicationStatus";

interface Props {
  event: ApplicationEvent;
  footer?: React.ReactNode; // сюда передаем либо статус, либо кнопки
}

const ApplicationEventCard: React.FC<Props> = ({ event, footer }) => {

    const startDate = new Date(event.start_date).toLocaleDateString();
    const endDate = new Date(event.end_date).toLocaleDateString();

    return (
        <Card className="d-flex mb-3">
            <Card.Body>
                <Col md={4}>
                    <Card.Img
                        // src можно передавать через event.image или что-то другое
                    />
                </Col>
                <Col md={8}>
                    <Card.Title>{event.name}</Card.Title>
                    <Card.Text>
                        <p><strong>Дата начала:</strong> {startDate}</p>
                        <p><strong>Дата окончания:</strong> {endDate}</p>
                    </Card.Text>
                </Col>
            </Card.Body>

            {/* Вариативная часть */}
            <div className="mt-3">
                {footer ? footer : (
                    <DefaultFooter event={event} />
                )}
            </div>
        </Card>
    )
}

// Можно вынести дефолтный футер
const DefaultFooter: React.FC<{ event: ApplicationEvent }> = ({ event }) => {
    const status = getStatusLabel(event.status);
    return (
        <div className="d-flex justify-content-between align-items-center">
            <small className="text-muted">
                *Правила проведения олимпиады смотрите*
                <a href="" className="text-decoration-underline">здесь</a>
            </small>
            <span className={`${status.className} fw-bold`}>
                {status.text}
            </span>
        </div>
    )
}

export default ApplicationEventCard;
