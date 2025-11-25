// pages/RegionalStagesPage.tsx
import { useEffect, useState } from "react";
import { Container, Spinner, Alert } from "react-bootstrap";
import { useNavigate } from 'react-router-dom';
import RegionalStageList from "./components/RegionalStageList";
import { useAuth } from "../../../Helpers/AuthContext";
import { fetchRegionalStages } from "../../../../requests/EventsRequests";

interface MyEvent {
    id?: string;
    name: string;
    start_date: string;
    end_date: string;
    previous_event_id?: string;
    event_type: string;
    subject?: number;
    class_number?: number;
    additional_info?: string;
    events?: MyEvent[];
}


const RegionalStagesPage: React.FC = () => {
  const [stages, setStages] = useState<MyEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();
  
  useEffect(() => {
  fetchRegionalStages()
    .then(setStages)
    .catch(err => setError((err as Error).message))
    .finally(() => setLoading(false));
}, []);


  const handleStageClick = (id: string) => {
  //console.log("Переход к этапу:", id);
  navigate(`/PersonalAccount/${id}`);
};

  return (
    <Container className="py-4">
      <h2 className="mb-5 text-center">Региональные этапы</h2>

      {loading && <Spinner animation="border" />}
      {error && <Alert variant="danger">{error}</Alert>}
      {!loading && !error && (
        <RegionalStageList stages={stages} onStageClick={handleStageClick} />
      )}
    </Container>
  );
};

export default RegionalStagesPage;