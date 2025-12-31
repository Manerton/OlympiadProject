import React, { useEffect, useState } from "react";
import { Table, Alert, Badge } from "react-bootstrap";
import axios from "axios";
import { APPLICATION, USER, SCHOOLS, API_CONFIG, AUTH } from "../../../config/api";
import { useAuth } from "../../Helpers/AuthContext";

// ===== Константы отображения =====
const CITIZENSHIP_TEXT: Record<number, string> = { 1: "Россия", 2: "Другое" };
const DISABILITY_TEXT: Record<number, string> = { 1: "Нет", 2: "Есть" };
const STATUS_TEXT: Record<number, string> = { 1: "Не обработано", 2: "Одобрено", 3: "Отклонено" };
const GENDER_TEXT: Record<number, string> = { 1: "Мужской", 2: "Женский" };

// ===== Типы =====
interface RawApplication {
    id: string;
    userId: string;
    schoolId: string;
    eventId: string;
    profile: string;
    class_participation: number;
    status: number;
    code: string;
    submittedAt: string;
}

interface District {
    id: string;
    name: string;
    region: number;
}

interface AggregatedApplication {
    id: string;
    firstname: string;
    surname: string;
    patronymic: string;
    email: string;
    phone: string;
    birthdate: string;
    gender: string;
    classNumber: number;
    citizenship: string;
    disability: string;
    schoolName: string;
    districtName: string;
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

async function axiosGetAllDistricts(
    token: string,
    regionId: number = 30
): Promise<Map<string, string>> {
    try {
        const res = await axios.get(`${AUTH.district}/${regionId}`, {
            headers: { Authorization: `Bearer ${token}` },
            withCredentials: true
        });

        const districts: District[] = res.data.data;
        const districtMap = new Map<string, string>();

        districts.forEach(d => districtMap.set(d.id, d.name));
        return districtMap;
    } catch (error) {
        console.error("Ошибка при получении районов:", error);
        return new Map();
    }
}

// ===== Компонент =====
const ApplicationsPage: React.FC = () => {
    const { accessToken } = useAuth();
    const [data, setData] = useState<AggregatedApplication[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (accessToken) fetchData();
    }, [accessToken]);

    const fetchData = async () => {
        if (!accessToken) {
            setError("Отсутствует токен авторизации");
            setLoading(false);
            return;
        }

        try {
            setLoading(true);
            setError(null);

            const districts = await axiosGetAllDistricts(accessToken);
            const raw = await axiosGetAllApplications(accessToken);

            const aggregated = await Promise.all(
                raw.map(async (app) => {
                    try {
                        const user = await axiosGetUser(accessToken, app.userId);
                        const school = await axiosGetSchool(accessToken, user.school_id);
                        const event = await axiosGetEvent(accessToken, app.eventId);

                        const districtName =
                            school.district_id && districts.has(school.district_id)
                                ? districts.get(school.district_id)!
                                : "Не указан";

                        return {
                            id: app.id,
                            surname: `${user.surname}`,
                            firstname: `${user.firstname}`,
                            patronymic: `${user.patronymic}`,
                            email: user.email ?? "—",
                            phone: user.phone ?? "—",
                            birthdate: user.birthdate,
                            gender: GENDER_TEXT[user.gender] ?? user.gender,
                            classNumber: user.class_number,
                            citizenship: CITIZENSHIP_TEXT[user.citizenship] ?? user.citizenship,
                            disability: DISABILITY_TEXT[user.disability] ?? user.disability,
                            schoolName: school.name,
                            districtName,
                            olympiadName: event.name,
                            profile: app.profile,
                            category: app.class_participation,
                            status: app.status,
                            code: app.code,
                            submittedAt: app.submittedAt
                        } as AggregatedApplication;
                    } catch {
                        return {
                            id: app.id,
                            surname: "Ошибка",
                            firstname: "Ошибка",
                            patronymic: "Ошибка",
                            email: "—",
                            phone: "—",
                            birthdate: "",
                            gender: "Ошибка",
                            classNumber: 0,
                            citizenship: "Ошибка",
                            disability: "Ошибка",
                            schoolName: "Ошибка",
                            districtName: "Ошибка",
                            olympiadName: "Ошибка",
                            profile: null,
                            category: app.class_participation,
                            status: app.status,
                            code: app.code,
                            submittedAt: app.submittedAt
                        } as AggregatedApplication;
                    }
                })
            );

            setData(aggregated);
        } catch (e: any) {
            setError(e.response?.data?.message || e.message || "Ошибка загрузки");
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
                    <th>Фамилия</th>
                    <th>Имя</th>
                    <th>Отчество</th>
                    <th>Email</th>
                    <th>Телефон</th>
                    <th>Школа</th>
                    <th>Муниципалитет</th>
                    <th>Олимпиада</th>
                    <th>Профиль</th>
                    <th>Класс обучения</th>
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
                        <td>{a.surname}</td>
                        <td>{a.firstname}</td>
                        <td>{a.patronymic}</td>
                        <td>{a.email}</td>
                        <td>{a.phone}</td>
                        <td>{a.schoolName}</td>
                        <td>{a.districtName}</td>
                        <td>{a.olympiadName}</td>
                        <td>{a.profile ?? "—"}</td>
                        <td className="text-center">
                            {a.classNumber}
                        </td>
                        <td className="text-center">
                           {a.category}
                        </td>
                        <td>{a.gender}</td>
                        <td>{a.citizenship}</td>
                        <td>{a.disability}</td>
                        <td>
                            <Badge>{STATUS_TEXT[a.status]}</Badge>
                        </td>
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
