import type React from "react"
import { taskTypes } from "../../../../../../dictionary/taskType"
import { useState } from "react"
import AppealForm from "./AppealForm"
import type { Task } from "../../../../../types/task"
import type { Appeal } from "../../../../../types/appeal"

type Props = {
    appeal: Appeal 
    onBack: () => void;
}

const AppealCreate: React.FC<Props> = ({appeal, onBack}) => {

    const handleSubmit = (e: React.FormEvent) => {

    }

    console.log(appeal)

    const [selectType, setSelectType] = useState<number | "">("");
    const [tasks, setTasks] = useState<Task[]>([])
    const [selectTask, setSelectTask] = useState<number | "">("");
    const [taskScore, setTaskScore] = useState<number>(0)
    const [reasonAppeal, setReasonAppeal] = useState<string>("")

    return (
        <div>
            <button className="btn btn-link mb-3" onClick={onBack}>
                Вернуться
            </button>
            <form onSubmit={handleSubmit} className="p-3">
                <AppealForm selectTask={selectTask} 
                tasks={""}
                setSelectTask={setSelectTask} 
                selectType={selectType} 
                setSelectType={setSelectType}
                taskScore={taskScore}
                reasonAppeal={reasonAppeal}
                setReasonAppeal={setReasonAppeal}
                />
                <button className="btn btn-primary" type="submit">Отправить</button>
            </form>
        </div>
       
    )

}

export default AppealCreate

