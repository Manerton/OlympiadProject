import React, { useEffect, useState } from "react";
import { Table, Alert, Badge, Button } from "react-bootstrap";
import axios from "axios";
import { API_CONFIG } from "../../../config/api";
import { useAuth } from "../../Helpers/AuthContext";

// ===== Константы отображения =====
const CITIZENSHIP_TEXT: Record<number, string> = { 1: "Россия", 2: "Другое" };
const DISABILITY_TEXT: Record<number, string> = { 1: "Нет", 2: "Есть" };
const STATUS_TEXT: Record<number, string> = { 
  1: "Не обработано", 
  2: "Одоброено", 
  3: "Отклонено" 
};
const GENDER_TEXT: Record<number, string> = { 1: "М", 2: "Ж" };

// ===== Типы =====
interface AggregatedApplicationResponse {
    id: string;
    firstname: string;
    surname: string;
    patronymic: string;
    email: string;
    phone: string;
    birthdate: string;
    gender: number;
    classNumber: number;
    citizenship: number;
    disability: number;
    schoolName: string;
    districtName: string;
    olympiadName: string;
    profile?: string | null;
    category: number;
    status: number;
    code: string;
    submittedAt: string;
}

interface ApiResponse<T> {
    status: string;
    status_code: number;
    data: T;
    error?: string;
}

// Вспомогательный тип для отображения
interface DisplayApplication {
    id: string;
    firstname: string;
    surname: string;
    patronymic: string;
    email: string;
    phone: string;
    birthdate: string;
    gender: string; // преобразованное
    classNumber: number;
    citizenship: string; // преобразованное
    disability: string; // преобразованное
    schoolName: string;
    districtName: string;
    olympiadName: string;
    profile?: string | null;
    category: number;
    status: number;
    statusText: string; // преобразованное
    code: string;
    submittedAt: string;
}

// ===== API helper =====
async function axiosGetAllAggregatedApplications(token: string): Promise<AggregatedApplicationResponse[]> {
    try {
        const res = await axios.get<ApiResponse<AggregatedApplicationResponse[]>>(
            API_CONFIG.ALLAPPLICATIONS,
            {
                headers: { Authorization: `Bearer ${token}` },
                withCredentials: true
            }
        );
        
        if (res.data.status_code !== 200) {
            throw new Error(res.data.error || `Ошибка ${res.data.status_code} при получении данных`);
        }
        
        return res.data.data || [];
    } catch (error: any) {
        console.error("API Error:", error);
        if (error.response?.data?.error) {
            throw new Error(error.response.data.error);
        }
        if (error.response?.data?.message) {
            throw new Error(error.response.data.message);
        }
        throw new Error("Не удалось загрузить данные. Проверьте подключение к интернету.");
    }
}

// ===== Вспомогательные функции =====
const formatDate = (dateString: string): string => {
    try {
        const date = new Date(dateString);
        return date.toLocaleDateString('ru-RU');
    } catch {
        return dateString;
    }
};

const getStatusBadgeVariant = (status: number): string => {
    switch(status) {
        case 1: return "warning";    // Не обработано
        case 2: return "success";    // Одобрено
        case 3: return "danger";     // Отклонено
        default: return "secondary";
    }
};

// Преобразование сырых данных в отображаемые
const transformDataForDisplay = (data: AggregatedApplicationResponse[]): DisplayApplication[] => {
    return data.map(item => ({
        ...item,
        gender: GENDER_TEXT[item.gender] || `Неизвестно (${item.gender})`,
        citizenship: CITIZENSHIP_TEXT[item.citizenship] || `Неизвестно (${item.citizenship})`,
        disability: DISABILITY_TEXT[item.disability] || `Неизвестно (${item.disability})`,
        statusText: STATUS_TEXT[item.status] || `Неизвестно (${item.status})`,
        birthdate: formatDate(item.birthdate),
        submittedAt: formatDate(item.submittedAt)
    }));
};

// ===== Компонент =====
const ApplicationsPage: React.FC = () => {
    const { accessToken } = useAuth();
    const [rawData, setRawData] = useState<AggregatedApplicationResponse[]>([]);
    const [displayData, setDisplayData] = useState<DisplayApplication[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (accessToken) {
            fetchData();
        } else {
            setError("Отсутствует токен авторизации");
            setLoading(false);
        }
    }, [accessToken]);

    const fetchData = async () => {
        if (!accessToken) return;

        try {
            setLoading(true);
            setError(null);

            const data = await axiosGetAllAggregatedApplications(accessToken);
            setRawData(data);
            setDisplayData(transformDataForDisplay(data));
        } catch (err: any) {
            console.error("Ошибка при загрузке данных:", err);
            setError(err.message || "Произошла ошибка при загрузке данных");
        } finally {
            setLoading(false);
        }
    };

    const refreshData = () => {
        if (accessToken) {
            fetchData();
        }
    };

    if (loading) {
        return (
            <div className="container py-5">
                <div className="text-center">
                    <div className="spinner-border text-primary" role="status">
                        <span className="visually-hidden">Загрузка...</span>
                    </div>
                    <p className="mt-3">Загрузка данных...</p>
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="container py-5">
                <Alert variant="danger">
                    <Alert.Heading>Ошибка загрузки данных</Alert.Heading>
                    <p>{error}</p>
                    <hr />
                    <div className="d-flex justify-content-end">
                        <Button onClick={refreshData} variant="outline-danger">
                            Попробовать снова
                        </Button>
                    </div>
                </Alert>
            </div>
        );
    }

    return (
        <div className="container py-4">
            <div className="d-flex justify-content-between align-items-center mb-4">
                <h3 className="mb-0">Все заявки</h3>
                <div className="d-flex gap-2">
                    <Button 
                        variant="outline-primary" 
                        onClick={refreshData}
                        size="sm"
                    >
                        Обновить
                    </Button>
                    <Badge bg="light" text="dark" className="fs-6">
                        Всего: {displayData.length}
                    </Badge>
                </div>
            </div>

            {displayData.length === 0 ? (
                <Alert variant="info">
                    Заявок не найдено
                </Alert>
            ) : (
                <>
                    <div className="table-responsive">
                        <Table bordered hover size="sm" className="align-middle">
                            <thead className="table-light">
                            <tr>
                                <th>№</th>
                                <th>Фамилия</th>
                                <th>Имя</th>
                                <th>Отчество</th>
                                <th>Дата рождения</th>
                                <th>Email</th>
                                <th>Телефон</th>
                                <th>Школа</th>
                                <th>Муниципалитет</th>
                                <th>Олимпиада</th>
                                <th>Профиль</th>
                                <th>Класс</th>
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
                            {displayData.map((app, index) => (
                                <tr key={app.id}>
                                    <td className="text-center fw-bold">{index + 1}</td>
                                    <td>{app.surname}</td>
                                    <td>{app.firstname}</td>
                                    <td>{app.patronymic}</td>
                                    <td>{app.birthdate}</td>
                                    <td>
                                        <small className="text-muted">{app.email}</small>
                                    </td>
                                    <td>
                                        <small>{app.phone}</small>
                                    </td>
                                    <td>
                                        <div className="text-truncate" style={{maxWidth: '200px'}} 
                                             title={app.schoolName}>
                                            {app.schoolName}
                                        </div>
                                    </td>
                                    <td>
                                        <small>{app.districtName}</small>
                                    </td>
                                    <td>
                                        <div className="text-truncate" style={{maxWidth: '150px'}}
                                             title={app.olympiadName}>
                                            {app.olympiadName}
                                        </div>
                                    </td>
                                    <td>
                                        {app.profile ? (
                                            <Badge bg="info" text="dark" className="fw-normal">
                                                {app.profile}
                                            </Badge>
                                        ) : (
                                            <span className="text-muted">—</span>
                                        )}
                                    </td>
                                    <td className="text-center">
                                        <Badge bg="secondary">{app.classNumber}</Badge>
                                    </td>
                                    <td className="text-center">
                                        <Badge bg="primary">{app.category}</Badge>
                                    </td>
                                    <td className="text-center">{app.gender}</td>
                                    <td>{app.citizenship}</td>
                                    <td>{app.disability}</td>
                                    <td>
                                        <Badge bg={getStatusBadgeVariant(app.status)}>
                                            {app.statusText}
                                        </Badge>
                                    </td>
                                    <td>
                                        {app.code ? (
                                            <code className="bg-light p-1 rounded">{app.code}</code>
                                        ) : (
                                            <span className="text-muted">—</span>
                                        )}
                                    </td>
                                    <td>
                                        <small className="text-muted">{app.submittedAt}</small>
                                    </td>
                                </tr>
                            ))}
                            </tbody>
                        </Table>
                    </div>
                    
                    <div className="d-flex justify-content-between align-items-center mt-3">
                        <div className="text-muted small">
                            Показано заявок: <strong>{displayData.length}</strong>
                        </div>
                        <Button 
                            variant="link" 
                            onClick={() => window.scrollTo({top: 0, behavior: 'smooth'})}
                            size="sm"
                        >
                            ↑ Наверх
                        </Button>
                    </div>
                </>
            )}
        </div>
    );
};

export default ApplicationsPage;