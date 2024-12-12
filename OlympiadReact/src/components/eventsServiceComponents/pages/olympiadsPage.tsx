import { useParams } from "react-router-dom";
import BaseEventPage from "./baseEventPage";
import { OLYMPIAD } from "../../../types/event";
import { RoleProvider } from "../../RoleContext";

function OlympiadsPage() {
    const { id } = useParams<{ id: string }>();
    return (
        <RoleProvider>
            <BaseEventPage type={OLYMPIAD} selectedEventId={id} pageName="Олимпиады" showSubjectField={true}/>

        </RoleProvider>
    )
}

export default OlympiadsPage