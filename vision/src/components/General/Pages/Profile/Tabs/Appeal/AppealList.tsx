import type React from "react";
import AppealCart from "./AppealCart";
import type { Appeal } from "../../../../../types/appeal";
import { taskTypeDict } from "../../../../../../dictionary/taskType";

const AppealList: React.FC = () => {

    const AppealData: Appeal[] = [
        {
            TaskType: taskTypeDict.Testing,
            TaskID: 101,
            Reason: "Ошибочный балл",
            Status: 0
        },
        {
            TaskType: taskTypeDict.Practic,
            TaskID: 202,
            Reason: "Недочет в оценке",
            Status: 1
        },
        {
            TaskType: taskTypeDict.Oral,
            TaskID: 303,
            Reason: "Не учтены дополнительные материалы",
            Status: 2
        }
    ];

    return (
        <div>
            {AppealData.map((appeal) => (
               <AppealCart key={appeal.TaskID} appeal={appeal}/>
            ))}
        </div>
    )
}

export default AppealList