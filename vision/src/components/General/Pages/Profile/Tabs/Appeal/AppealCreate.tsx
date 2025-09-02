import type React from "react";
import { useEffect, useState } from "react";
import AppealForm from "./AppealForm";
import type { Task } from "../../../../../types/task";
import { useAuth } from "../../../../../Helpers/AuthContext";
import { axiosCreateAppeal, axiosResultGetByEventUser } from "../../../../../../requests/ResultRequests";
import { data, useParams } from "react-router-dom";
import type { Result } from "../../../../../types/result";
import type { CreateAppealRequest } from "../../../../../types/appeal";


const AppealCreate: React.FC = () => {

    const { accessToken, user } = useAuth();


    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        const appealData: CreateAppealRequest = {
            reason: reasonAppeal,
            taskId: selectTask as number,
            userId: user?.id as string,
        };

        try {
            await axiosCreateAppeal(accessToken!, appealData);
            alert("Апелляция отправлена");

            // сброс формы
            setSelectType("");
            setSelectTask("");
            setTaskScore(0);
            setReasonAppeal("");
        } catch (err) {
            console.error("Ошибка отправки апелляции", err);
        } finally {
            console.log("Отправленные данные:", appealData);
        }
    };


    const [allTasks, setAllTasks] = useState<Result[]>([]);
    const [tasks, setTasks] = useState<Task[]>([]);
    const [selectType, setSelectType] = useState<number | "">("");
    const [selectTask, setSelectTask] = useState<number | "">("");
    const [taskScore, setTaskScore] = useState<number>(0);
    const [reasonAppeal, setReasonAppeal] = useState<string>("");

    const { eventId } = useParams();

    useEffect(() => {
        async function fetchResult() {
            try {
                const result = await axiosResultGetByEventUser(accessToken!, eventId!, user?.id!);
                setAllTasks(result);
            } catch (err) {
                console.error("Ошибка загрузки задач", err);
            }
        }
        fetchResult();
    }, [accessToken, user, eventId]);

    // когда выбран тип — фильтруем список заданий
    useEffect(() => {
        if (selectType !== "") {
            const filtered = allTasks
                .filter(t => t.taskType === selectType)
                .map(t => ({ id: t.taskId, name: `Задание ${t.taskName}` }));
            setTasks(filtered);
            setSelectTask(""); // сбрасываем выбранное задание
            setTaskScore(0);   // сбрасываем балл
        }
    }, [selectType, allTasks]);

    // когда выбрано задание — подставляем баллы
    useEffect(() => {
        if (selectTask !== "") {
            const found = allTasks.find(t => t.taskId === selectTask);
            if (found) setTaskScore(found.points);
        }
    }, [selectTask, allTasks]);

    return (
        <div>
            <form onSubmit={handleSubmit} className="p-3">
                <AppealForm
                    selectTask={selectTask}
                    tasks={tasks}
                    setSelectTask={setSelectTask}
                    selectType={selectType}
                    setSelectType={setSelectType}
                    taskScore={taskScore}
                    reasonAppeal={reasonAppeal}
                    setReasonAppeal={setReasonAppeal}
                />
                <button className="btn btn-primary" type="submit">
                    Отправить
                </button>
            </form>
        </div>
    );
};

export default AppealCreate;
