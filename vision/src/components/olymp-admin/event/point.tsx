import React, { useEffect, useState } from "react";
import axios from "axios";
import { useNavigate, useParams } from "react-router-dom";

interface Task {
    id: string;
    number: string;
    max_points: number;
}

interface TaskAttendance {
    id: string;
    points: string | null;
    task: Task;
}

interface Person {
    fullFio: string;
}

interface Application {
    code: string;
}

interface TableRow {
    person: Person;
    application: Application;
    taskAttendances: TaskAttendance[];
}

const EventPoint: React.FC = () => {
    const { eventId } = useParams<{ eventId: string }>();
    const [tasks, setTasks] = useState<Task[]>([]);
    const [table, setTable] = useState<TableRow[]>([]);
    const [loading, setLoading] = useState<boolean>(true);
    const navigate = useNavigate();

    const token =
        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQG1haWwucnUiLCJleHAiOjE3ODU0OTE5MzksImlkIjoiMGU2OTkxOTQtZjc4MS00NWE2LTg3Y2YtNTRhOTYyMzI1Y2YyIiwicm9sZSI6MX0.-bc6ZKSP6Lbv6rYO89ZV65iWVHxCrFlUDPjM81N1Dyc";

    const fetchData = async () => {
        try {
            setLoading(true);
            const response = await axios.get(
                `http://olymp-admin-v2/api/event/points/${eventId}`,
                {
                    headers: {
                        Authorization: token,
                    },
                    withCredentials: true,
                }
            );

            setTasks(response.data.tasks || []);
            setTable(response.data.table || []);
        } catch (error) {
            console.error("Error fetching event points:", error);
        } finally {
            setLoading(false);
        }
    };

    const handlePointsChange = async (
        taskAttendanceId: string,
        points: string
    ) => {
        const cleanedPoints = points.replace(/[^0-9]/g, "");
        try {
            await axios.post(
                "http://olymp-admin-v2/api/event/change-score",
                {
                    task_attendance_id: taskAttendanceId,
                    points: cleanedPoints,
                },
                {
                    headers: {
                        Authorization: token,
                        "Content-Type": "application/json",
                    },
                    withCredentials: true,
                }
            );
        } catch (error) {
            console.error("Error updating points:", error);
        }
    };

    useEffect(() => {
        fetchData();
    }, [eventId]);

    if (loading) {
        return <div className="text-center mt-4">Загрузка...</div>;
    }

    return (
        <div className="container mt-4">
            <button
                className="btn btn-sm btn-primary mb-3"
                onClick={() => navigate(`/olymp-admin/event/show/${eventId}`)}
            >
                Перейти в карточку олимпиады
            </button>

            <h1>Выставление баллов</h1>

            <table className="table table-bordered table-striped mt-3">
                <thead>
                    <tr>
                        <th>ФИО участника</th>
                        <th>Код участника</th>
                        {tasks.map((task) => (
                            <th key={task.id}>{task.number}</th>
                        ))}
                    </tr>
                </thead>
                <tbody>
                    {table.map((row, rowIndex) => (
                        <tr key={rowIndex}>
                            <td>{row.person.fullFio}</td>
                            <td>{row.application.code}</td>
                            {row.taskAttendances.map((ta) => (
                                <td key={ta.id}>
                                    <input
                                        type="text"
                                        value={ta.points || ""}
                                        maxLength={ta.task.max_points.toString().length}
                                        onChange={(e) =>
                                            handlePointsChange(
                                                ta.id,
                                                e.target.value
                                            )
                                        }
                                        onInput={(e) => {
                                            const target =
                                                e.target as HTMLInputElement;
                                            target.value = target.value.replace(
                                                /[^0-9]/g,
                                                ""
                                            );
                                        }}
                                        className="form-control form-control-sm"
                                    />
                                </td>
                            ))}
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
};

export default EventPoint;