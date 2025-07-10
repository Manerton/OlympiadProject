import { useParams } from "react-router-dom";
import BaseEventPage from "./baseEventPage";
import { OLYMPIAD } from "../../../types/event";

function OlympiadsPage() {
    const { id } = useParams<{ id: string }>();
    return (
        <BaseEventPage EventType={OLYMPIAD} selectedEventId={id} pageName="Олимпиады по предметам" showSubjectField={true} />
    )
}

export default OlympiadsPage