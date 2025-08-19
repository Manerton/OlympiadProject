import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate, useParams } from 'react-router-dom';

interface AttendanceItem {
    person: {
        firstname?: string;
        surname?: string;
        patronymic?: string;
    };
    attendance: {
        id: number;
        status: string | number;
    };
}

interface AttendanceStatuses {
    [key: string]: string;
}

interface Event {
    id: string;
    name: string;
}

const EventAttendance: React.FC = () => {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const [event, setEvent] = useState<Event | null>(null);
    const [attendanceData, setAttendanceData] = useState<AttendanceItem[]>([]);
    const [attendanceStatuses, setAttendanceStatuses] = useState<AttendanceStatuses>({});
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQG1haWwucnUiLCJleHAiOjE3ODU0OTE5MzksImlkIjoiMGU2OTkxOTQtZjc4MS00NWE2LTg3Y2YtNTRhOTYyMzI1Y2YyIiwicm9sZSI6MX0.-bc6ZKSP6Lbv6rYO89ZV65iWVHxCrFlUDPjM81N1Dyc';
        const fetchData = async () => {
            try {
                const response = await axios.get(`http://localhost:8080/api/event/attendance/${id}`, {
                    headers: { 'Authorization': token }
                });

                setEvent(response.data.event);
                
                setAttendanceData(response.data.data || []);
                setAttendanceStatuses(response.data.attendanceStatuses || {});
                setLoading(false);
                
            } catch (error) {
                console.error('Error fetching attendance data:', error);
                setLoading(false);
            }
        };

        fetchData();
    }, [id]);

    const handleStatusChange = async (attendanceId: number, newStatus: string) => {
        const token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQG1haWwucnUiLCJleHAiOjE3ODU0OTE5MzksImlkIjoiMGU2OTkxOTQtZjc4MS00NWE2LTg3Y2YtNTRhOTYyMzI1Y2YyIiwicm9sZSI6MX0.-bc6ZKSP6Lbv6rYO89ZV65iWVHxCrFlUDPjM81N1Dyc';
        
        try {
            await axios.post('http://localhost:8080/api/event/change-attendance', {
                attendance_id: attendanceId,
                status: newStatus,
                eventId: id
            }, {
                headers: {
                    'Authorization': token,
                    'Content-Type': 'application/json'
                },
                withCredentials: true
            });
            
            // Обновляем только статус внутри attendance
            setAttendanceData(prevData => 
                prevData.map(item => 
                    item.attendance.id === attendanceId 
                        ? { ...item, attendance: { ...item.attendance, status: newStatus } }
                        : item
                )
            );
        } catch (error) {
            console.error('Error updating attendance status:', error);
            alert('Не удалось обновить статус явки');
        }
    };
    {attendanceData.map((item, index) => (
            console.log(item)
    ))};
    const getFullName = (person: { firstname?: string; surname?: string; patronymic?: string }) => {
        return `${person.firstname || ''} ${person.surname || ''} ${person.patronymic || ''}`.trim();
    };

    if (loading) return <div className="text-center mt-4"><p>Загрузка...</p></div>;
    if (!event) return <div className="text-center mt-4"><p>Олимпиада не найдена</p></div>;

    return (
        <div className="event-attendance container mt-4">
            <button 
                onClick={() => navigate(`/olymp-admin/event/show/${event.id}`)}
                className="btn btn-sm btn-primary mb-3"
            >
                Перейти в карточку олимпиады
            </button>

            <table className="table table-bordered table-striped">
                <thead>
                    <tr>
                        <th>#</th>
                        <th>ФИО участника</th>
                        <th>Явка</th>
                    </tr>
                </thead>
                <tbody>
                    {attendanceData.map((item, index) => (
                        <tr key={item.attendance.id}>
                            <td>{index + 1}</td>
                            <td>{getFullName(item.person)}</td>
                            <td>
                                <select
                                    className="form-control"
                                    value={item.attendance.status}
                                    onChange={(e) => handleStatusChange(item.attendance.id, e.target.value)}
                                >
                                    {Object.entries(attendanceStatuses).map(([key, status]) => (
                                        <option key={key} value={key}>
                                            {status}
                                        </option>
                                    ))}
                                </select>
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
};

export default EventAttendance;
