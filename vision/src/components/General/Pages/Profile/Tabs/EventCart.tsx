import { Card, Col, Row } from "react-bootstrap";
import type { ApplicationEvent } from "../../../../types/event"
import CardImage from "../../Events/components/CardImage";
import { StatusIcon } from "../../../../Helpers/StatusBlock";

interface Props {
    event: ApplicationEvent;
    footer?: React.ReactNode; // сюда передаем либо статус, либо кнопки
}

const EventCard: React.FC<Props> = ({ event, footer }) => {

    const startDate = new Date(event.start_date).toLocaleDateString();
    const endDate = new Date(event.end_date).toLocaleDateString();


    return (
        <Card className="mb-3" >
            <Card.Body>
                <Row className="">
                    <Col md={3}>
                        <CardImage subjectId={event.subject ?? 0} />
                    </Col>
                    <Col md={7}>
                        <Card.Title>{event.name}</Card.Title>
                        <Card.Text>
                            <p>
                                <strong>Дата начала:</strong> {startDate}
                            </p>
                            <p>
                                <strong>Дата окончания:</strong> {endDate}
                            </p>
                        </Card.Text>
                    </Col>
                    <Col md={2} className="">
                        {footer ? footer : <DefaultFooter event={event} />}
                    </Col>
                </Row>
            </Card.Body>
        </Card>
    )
}

// Можно вынести дефолтный футер
const DefaultFooter: React.FC<{ event: ApplicationEvent }> = ({ event }) => {
    return (
        <div className="align-items-center">
            <small className="text-muted">
                *Правила проведения олимпиады смотрите*
                <a href="" className="text-decoration-underline">здесь</a>
            </small>
            <StatusIcon status={event.status} />
        </div>
    )
}

export default EventCard;
