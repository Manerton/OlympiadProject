import React, { useEffect, useState } from "react";
import { Table, Alert, Badge } from "react-bootstrap";
import axios from "axios";
import { APPLICATION, USER, SCHOOLS, API_CONFIG } from "../../../config/api";
import {useAuth} from "../../Helpers/AuthContext";

// ===== Константы отображения =====
const CITIZENSHIP_TEXT: Record<number, string> = { 1: "Россия", 2: "Другое" };
const DISABILITY_TEXT: Record<number, string> = { 0: "Нет", 1: "Есть" };
const STATUS_TEXT: Record<number, string> = { 1: "Не обработано", 2: "Одобрено", 3: "Отклонено" };
const GENDER_TEXT: Record<number, string> = { 1: "Мужской", 2: "Женский" };

// ===== Типы =====
interface RawApplication {
    id: string;
    userId: string;
    schoolId: string;
    eventId: string;
    class_participation: number;
    status: number;
    code: string;
    submittedAt: string;
}

interface AggregatedApplication {
    id: string;
    fio: string;
    birthdate: string;
    gender: string;
    classNumber: number;
    citizenship: string;
    disability: string;
    schoolName: string;
    olympiadName: string;
    profile?: string | null;
    category: number;
    status: number;
    code: string;
    submittedAt: string;
}

// ===== API helpers =====
async function axiosGetAllApplications(token: string): Promise<RawApplication[]> {
    const res = await axios.get(APPLICATION.getALL, {
        headers: { Authorization: `Bearer ${token}` },
        withCredentials: true
    });
    return res.data.data as RawApplication[];
}

async function axiosGetUser(token: string, id: string) {
    const res = await axios.get(`${USER.infoParticipant}/${id}`, {
        headers: { Authorization: `Bearer ${token}` },
        withCredentials: true
    });
    return res.data.data;
}

async function axiosGetSchool(token: string, id: string) {
    const res = await axios.get(`${SCHOOLS.byId}/${id}`, {
        headers: { Authorization: `Bearer ${token}` },
        withCredentials: true
    });
    return res.data.data;
}

async function axiosGetEvent(token: string, id: string) {
    const res = await axios.get(`${API_CONFIG.EVENT}/${id}`, {
        headers: { Authorization: `Bearer ${token}` },
        withCredentials: true
    });
    return res.data.data;
}

// ===== Компонент =====
const ApplicationsPage: React.FC = () => {
    const { accessToken } = useAuth();
    const [data, setData] = useState<AggregatedApplication[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        fetchData();
    }, []);

    const fetchData = async () => {
        try {
            setLoading(true);
            const raw = await axiosGetAllApplications(accessToken!);

            const aggregated = await Promise.all(
                raw.map(async (app) => {
                    const user = await axiosGetUser(accessToken!, app.userId);
                    const school = await axiosGetSchool(accessToken!, user.school_id);
                    const event = await axiosGetEvent(accessToken!, app.eventId);

                    return {
                        id: app.id,
                        fio: `${user.surname} ${user.firstname} ${user.patronymic}`,
                        birthdate: user.birthdate,
                        gender: GENDER_TEXT[user.gender] ?? user.gender,
                        classNumber: user.class_number,
                        citizenship: CITIZENSHIP_TEXT[user.citizenship] ?? user.citizenship,
                        disability: DISABILITY_TEXT[user.disability] ?? user.disability,
                        schoolName: school.name,
                        olympiadName: event.name,
                        profile: event.profiles,
                        category: app.class_participation,
                        status: app.status,
                        code: app.code,
                        submittedAt: app.submittedAt
                    } as AggregatedApplication;
                })
            );

            setData(aggregated);
        } catch (e: any) {
            setError(e.message);
        } finally {
            setLoading(false);
        }
    };

    if (loading) return <div className="p-5 text-center">Загрузка...</div>;
    if (error) return <Alert variant="danger">{error}</Alert>;

    return (
        <div className="container py-4">
            <h3 className="mb-3">Все заявки</h3>
            <Table bordered hover responsive>
                <thead>
                <tr>
                    <th>№</th>
                    <th>ФИО</th>
                    <th>Школа</th>
                    <th>Олимпиада</th>
                    <th>Профиль</th>
                    <th>Класс участия</th>
                    <th>Пол</th>
                    <th>Гражданство</th>
                    <th>ОВЗ</th>
                    <th>Статус</th>
                    <th>Код</th>
                    <th>Дата подачи</th>
                </tr>
                </thead>
                <tbody>
                {data.map((a, i) => (
                    <tr key={a.id}>
                        <td>{i + 1}</td>
                        <td>{a.fio}</td>
                        <td>{a.schoolName}</td>
                        <td>{a.olympiadName}</td>
                        <td>{a.profile ?? "—"}</td>
                        <td className="text-center"><Badge bg="primary">{a.category}</Badge></td>
                        <td>{a.gender}</td>
                        <td>{a.citizenship}</td>
                        <td>{a.disability}</td>
                        <td><Badge>{STATUS_TEXT[a.status]}</Badge></td>
                        <td>{a.code}</td>
                        <td>{new Date(a.submittedAt).toLocaleDateString()}</td>
                    </tr>
                ))}
                </tbody>
            </Table>

            <div className="mt-3 text-muted small">
                Всего заявок: {data.length}
            </div>
        </div>
    );
};

export default ApplicationsPage;