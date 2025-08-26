// pages/RegionalStagesPage.tsx
import React, { useEffect, useState } from "react";
import { Container, Spinner, Alert } from "react-bootstrap";
import { useNavigate } from 'react-router-dom';
import RegionalStageList from "./components/RegionalStageList";
import axios from 'axios';
import { useAuth } from "../../../Helpers/AuthContext";
import {API_CONFIG} from "../../../../config/api.ts";
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
  const { accessToken } = useAuth();
  const navigate = useNavigate();
  useEffect(() => {
    const fetchStages = async () => {
      try {
        const res = await axios.get(API_CONFIG.REGIONAL, {
          headers: {
              'Authorization': `Bearer ${accessToken}`  // <- добавляем "Bearer " перед токеном
          },
          withCredentials: true
      });
        setStages(res.data.data);
      } catch (err) {
        setError((err as Error).message);
      } finally {
        setLoading(false);
      }
    };

    fetchStages();
  }, []);

  const handleStageClick = (id: string) => {
  console.log("Переход к этапу:", id);
  navigate(`/OlympiadsPage/${id}`);
};

  return (
    <Container className="py-4">
      <h2 className="mb-4">Региональные этапы</h2>

      {loading && <Spinner animation="border" />}
      {error && <Alert variant="danger">{error}</Alert>}
      {!loading && !error && (
        <RegionalStageList stages={stages} onStageClick={handleStageClick} />
      )}
    </Container>
  );
};

export default RegionalStagesPage;