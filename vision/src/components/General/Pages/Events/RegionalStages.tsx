// pages/RegionalStagesPage.tsx
import React, { useEffect, useState } from "react";
import { Container, Spinner, Alert } from "react-bootstrap";
import { useNavigate } from 'react-router-dom';
import RegionalStageList from "./components/RegionalStageList";
import axios from 'axios';

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


const API_URL = "http://172.16.1.39:8080/api/events/regional-stage";

const RegionalStagesPage: React.FC = () => {
  const [stages, setStages] = useState<MyEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQG1haWwucnUiLCJleHAiOjE3ODU0OTE5MzksImlkIjoiMGU2OTkxOTQtZjc4MS00NWE2LTg3Y2YtNTRhOTYyMzI1Y2YyIiwicm9sZSI6MX0.-bc6ZKSP6Lbv6rYO89ZV65iWVHxCrFlUDPjM81N1Dyc';
  const navigate = useNavigate();
  useEffect(() => {
    const fetchStages = async () => {
      try {
        const res = await axios.get(API_URL);
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