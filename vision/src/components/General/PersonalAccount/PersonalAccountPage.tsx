import { useState } from "react";
import { Collapse, Button, Card } from "react-bootstrap";

// Твои компоненты
import PersonalInfo from "./PersonalInfo";
import ApplicationEventPage from "./ApplicationEvent";
import OlympiadsSimpleTable from "./Olympiads";
import { eachWeekOfInterval } from "date-fns";

const ParticipantDashboard: React.FC = () => {
    const [openInfo, setOpenInfo] = useState(true);
    const [openApps, setOpenApps] = useState(false);
    const [openOlympiads, setOpenOlympiads] = useState(false);

    const [reloadFlag, setReloadFlag] = useState(0);

    const reloadAll = () => setReloadFlag((prev) => prev + 1);

    return (
        <div className="container mt-4">

            {/* ====== ЛИЧНАЯ ИНФОРМАЦИЯ ====== */}
            <Card className="mb-3 shadow-sm">
                <Card.Header
                    className="d-flex justify-content-between align-items-center"
                    style={{ cursor: "pointer" }}
                    onClick={() => setOpenInfo(!openInfo)}
                >
                    <h5 className="mb-0">Личная информация</h5>
                    <Button variant="link" className="p-0">
                        {openInfo ? "▲" : "▼"}
                    </Button>
                </Card.Header>

                <Collapse in={openInfo}>
                    <div>
                        <Card.Body>
                            <PersonalInfo />
                        </Card.Body>
                    </div>
                </Collapse>
            </Card>

            {/* ====== РАНЕЕ ПОДАННЫЕ ЗАЯВКИ ====== */}
            <Card className="mb-3 shadow-sm">
                <Card.Header
                    className="d-flex justify-content-between align-items-center"
                    style={{ cursor: "pointer" }}
                    onClick={() => setOpenApps(!openApps)}
                >
                    <h5 className="mb-0">Мои заявки</h5>
                    <Button variant="link" className="p-0">
                        {openApps ? "▲" : "▼"}
                    </Button>
                </Card.Header>

                <Collapse in={openApps}>
                    <div>
                        <Card.Body>
                            <ApplicationEventPage onApplied={reloadAll} reloadFlag={reloadFlag} />
                        </Card.Body>
                    </div>
                </Collapse>
            </Card>

            {/* ====== ОЛИМПИАДЫ ДЛЯ УЧАСТИЯ ====== */}
            <Card className="mb-3 shadow-sm">
                <Card.Header
                    className="d-flex justify-content-between align-items-center"
                    style={{ cursor: "pointer" }}
                    onClick={() => setOpenOlympiads(!openOlympiads)}
                >
                    <h5 className="mb-0">Подать новую заявку</h5>
                    <Button variant="link" className="p-0">
                        {openOlympiads ? "▲" : "▼"}
                    </Button>
                </Card.Header>

                <Collapse in={openOlympiads}>
                    <div>
                        <Card.Body>
                            <OlympiadsSimpleTable onApplied={reloadAll} reloadFlag={reloadFlag} />
                        </Card.Body>
                    </div>
                </Collapse>
            </Card>

        </div>
    );
};

export default ParticipantDashboard;
