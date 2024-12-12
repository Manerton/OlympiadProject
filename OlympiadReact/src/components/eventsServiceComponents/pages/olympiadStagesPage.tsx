import { useParams } from "react-router-dom";
import BaseEventPage from "./baseEventPage";
import { STAGE } from "../../../types/event";
import { RoleProvider } from "../../RoleContext";


function OlympiadStagesPage() {
    const { id } = useParams<{ id: string }>();
    return (
        <RoleProvider>
            <BaseEventPage type={STAGE} selectedEventId={id} pageName="Этапы олимпиады" />
        </RoleProvider>

    )
}

export default OlympiadStagesPage