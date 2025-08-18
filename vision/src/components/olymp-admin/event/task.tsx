import React, { useEffect, useState } from "react";
import axios from "axios";
import { useParams, useNavigate } from "react-router-dom";

interface Task {
    id: string;
    number: string;
    max_points: string;
}
interface Event {
    id: string;
    name: string;
}

interface NewTask {
    number: string;
    point: string;
}

const EventTask: React.FC = () => {
    const { id } = useParams<{ id: string }>();
    const [event, setEvent] = useState<Event | null>(null);
    const navigate = useNavigate();

    const [tasks, setTasks] = useState<Task[]>([]);
    const [newTasks, setNewTasks] = useState<NewTask[]>([
        { number: "", point: "" },
    ]);
    const [loading, setLoading] = useState<boolean>(true);
    console.log(id);
    const token =
        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQG1haWwucnUiLCJleHAiOjE3ODU0OTE5MzksImlkIjoiMGU2OTkxOTQtZjc4MS00NWE2LTg3Y2YtNTRhOTYyMzI1Y2YyIiwicm9sZSI6MX0.-bc6ZKSP6Lbv6rYO89ZV65iWVHxCrFlUDPjM81N1Dyc";
    const fetchTasks = async () => {
        try {
            setLoading(true);
            const response = await axios.get(
                `http://localhost:8080/api/event/task/${id}`,
                {
                    headers: { Authorization: token },
                    withCredentials: true,
                }
            );
            setEvent(response.data.event);
            
            setTasks(response.data.tasks || []);
        } catch (error) {
            console.error("Error fetching tasks:", error);
        } finally {
            setLoading(false);
        }
    };
    const handleAddField = () => {
        setNewTasks([...newTasks, { number: "", point: "" }]);
    };

    const handleRemoveField = (index: number) => {
        if (newTasks.length > 1) {
            setNewTasks(newTasks.filter((_, i) => i !== index));
        }
    };

    const handleChange = (
        index: number,
        field: keyof NewTask,
        value: string
    ) => {
        const updatedTasks = [...newTasks];
        updatedTasks[index][field] = value.replace(/[^0-9]/g, "");
        setNewTasks(updatedTasks);
    };

    const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
        // Преобразуем данные в нужный формат
        const requestData = {
            number: newTasks.map(task => task.number),
            point: newTasks.map(task => task.point) // Обратите внимание: у вас в интерфейсе поле называется 'point', но здесь опечатка 'point'
        };

        await axios.post(
            `http://localhost:8080/api/event/add-task/${id}`,
            requestData,
            {
                headers: {
                    Authorization: token,
                    "Content-Type": "application/json",
                },
                withCredentials: true,
            }
        );
        setNewTasks([{ number: "", point: "" }]);
        fetchTasks();
    } catch (error) {
        console.error("Error adding tasks:", error);
    }
};

    const handleDelete = async (taskId: string) => {
        if (!window.confirm("Вы уверены, что хотите удалить этот элемент?")) {
            return;
        }
        try {
            await axios.delete(
                `http://localhost:8080/api/event/delete-task/${taskId}`,
                {
                    headers: { Authorization: token },
                    withCredentials: true,
                }
            );
            fetchTasks();
        } catch (error) {
            console.error("Error deleting task:", error);
        }
    };

    useEffect(() => {
        fetchTasks();
    }, [id]);

    if (loading) {
        return <div className="text-center mt-4">Загрузка...</div>;
    }

    return (
        <div className="container mt-4">
            <button
                className="btn btn-sm btn-primary mb-3"
                onClick={() => navigate(`/olymp-admin/event/show/${id}`)}
            >
                Перейти в карточку олимпиады
            </button>

            <h1>Список заданий</h1>

            {/* Форма добавления */}
            <form onSubmit={handleSubmit} className="mb-4">
                {newTasks.map((task, index) => (
                    <div
                        className="row align-items-end mb-2"
                        key={`new-task-${index}`}
                    >
                        <div className="col-md-5">
                            <label className="form-label">
                                Номер задания
                            </label>
                            <input
                                type="text"
                                className="form-control"
                                value={task.number}
                                onChange={(e) =>
                                    handleChange(index, "number", e.target.value)
                                }
                                required
                            />
                        </div>
                        <div className="col-md-5">
                            <label className="form-label">
                                Количество баллов
                            </label>
                            <input
                                type="text"
                                className="form-control"
                                value={task.point}
                                onChange={(e) =>
                                    handleChange(index, "point", e.target.value)
                                }
                                required
                            />
                        </div>
                        <div className="col-md-2">
                            <button
                                type="button"
                                className="btn btn-danger"
                                onClick={() => handleRemoveField(index)}
                            >
                                −
                            </button>
                        </div>
                    </div>
                ))}
                <button
                    type="button"
                    id="add-field"
                    className="btn btn-primary mb-3"
                    onClick={handleAddField}
                >
                    +
                </button>
                <br />
                <button type="submit" className="btn btn-success">
                    Сохранить
                </button>
            </form>

            {/* Таблица существующих заданий */}
            <table className="table table-bordered table-striped">
                <thead>
                    <tr>
                        <th>#</th>
                        <th>Номер задания</th>
                        <th>Максимальное количество баллов</th>
                        <th>Действия</th>
                    </tr>
                </thead>
                <tbody>
                    {tasks.map((task, index) => (
                        <tr key={task.id}>
                            <td>{index + 1}</td>
                            <td>{task.number}</td>
                            <td>{task.max_points}</td>
                            <td>
                                <button
                                    className="btn btn-sm btn-danger"
                                    onClick={() => handleDelete(task.id)}
                                >
                                    Удалить
                                </button>
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
};

export default EventTask;