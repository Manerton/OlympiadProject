import React, { useEffect, useState } from 'react';
import { Container, Spinner, Alert, Form, Button, Pagination } from 'react-bootstrap';
import DatePicker from "react-datepicker";
import "react-datepicker/dist/react-datepicker.css";
import OlympiadList from './components/OlympiadList';
import axios from 'axios';
import { useNavigate,useParams } from 'react-router-dom';
import type { MyEvent } from '../../../types/event.ts';

const API_BASE = 'http://172.16.1.39:8080/api/events/child';
const TOKEN = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQG1haWwucnUiLCJleHAiOjE3ODU0OTE5MzksImlkIjoiMGU2OTkxOTQtZjc4MS00NWE2LTg3Y2YtNTRhOTYyMzI1Y2YyIiwicm9sZSI6MX0.-bc6ZKSP6Lbv6rYO89ZV65iWVHxCrFlUDPjM81N1Dyc';

const OlympiadsPage: React.FC = () => {
  const [olympiads, setOlympiads] = useState<MyEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [page, setPage] = useState(0);
  const [limit] = useState(6);
  const [order] = useState('name DESC');
  const [search, setSearch] = useState('');
  const [selectedDate, setSelectedDate] = useState<Date | null>(null);
  const [total, setTotal] = useState(0);
  const { id } = useParams();
  const navigate = useNavigate();

  useEffect(() => {
    const fetchOlympiads = async () => {
      setLoading(true);
      try {
        let url = `${API_BASE}/${id}?page=${page}&limit=${limit}&order=${order}`;
        if (search) {
          url += `&name=${encodeURIComponent(search)}`;
        }
        if (selectedDate) {
          const dateStr = selectedDate.toISOString().split('T')[0]; // формат YYYY-MM-DD
          url += `&date=${dateStr}`; // предполагается, что API принимает параметр date
        }
        const res = await axios.get(url, {
          headers: { Authorization: `Bearer ${TOKEN}` },
        });
        setOlympiads(res.data.data);
        setTotal(res.data.metadata);
      } catch (err) {
        setError((err as Error).message);
      } finally {
        setLoading(false);
      }
    };

    fetchOlympiads();
  }, [page, limit, order, search, selectedDate]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setPage(1);
  };

  const handleOlympiadClick = (id: string) => {
    navigate(`/stages/${id}`);
  };

 const totalPages = Math.ceil(total / limit);
//console.log('Всего страниц:', totalPages, 'Всего записей:', total, 'Лимит на страницу:', limit);


  return (
    <Container className="py-4">
      <h2 className="mb-4">СПИСОК ОЛИМПИАД</h2>

      <Form onSubmit={handleSearch} className="mb-4">
        <div className="d-flex align-items-center">
          <Form.Control
            type="text"
            placeholder="Название олимпиады"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />

          <div className="ms-2">
            <DatePicker
              selected={selectedDate}
              onChange={(date) => setSelectedDate(date)}
              dateFormat="dd.MM.yyyy"
              className="form-control"
              placeholderText="Выбрать дату"
              isClearable
            />
          </div>

          <Button variant="primary" type="submit" className="ms-2">
            Поиск
          </Button>
        </div>
      </Form>

      {loading && <Spinner animation="border" />}
      {error && <Alert variant="danger">{error}</Alert>}
      {!loading && !error && (
        <OlympiadList olympiads={olympiads} onOlympiadClick={handleOlympiadClick} />
      )}

      {totalPages > 1 && (
        <Pagination className="justify-content-center mt-4">
          <Pagination.Prev
            onClick={() => setPage((p) => Math.max(1, p - 1))}
            disabled={page === 1}
          />
          {Array.from({ length: totalPages }).map((_, index) => (
            <Pagination.Item
              key={index + 1}
              active={index + 1 === page}
              onClick={() => setPage(index + 1)}
            >
              {index + 1}
            </Pagination.Item>
          ))}
          <Pagination.Next
            onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
            disabled={page === totalPages}
          />
        </Pagination>
      )}
    </Container>
  );
};

export default OlympiadsPage;
