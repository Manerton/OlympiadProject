import { useParams } from "react-router-dom";
import BaseEventPage from "./baseEventPage";
import { CLASS } from "../../../types/event";


function OlympiadClassPage() {
    const { id } = useParams<{ id: string }>();
    return (
        <BaseEventPage EventType={CLASS} selectedEventId={id} pageName="Классы по предмету" showClassNumber={true} />
    )
}

export default OlympiadClassPage