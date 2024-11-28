import { useParams } from "react-router-dom";
import BaseEventPage from "./baseEventPage";
import { STAGE } from "../../../types/event";


function OlympiadStagesPage() {
    const { id } = useParams<{ id: string }>();
    return (
        <BaseEventPage type={STAGE} selectedEventId={id} pageName="Этапы олимпиады"/>
    )
}

export default OlympiadStagesPage