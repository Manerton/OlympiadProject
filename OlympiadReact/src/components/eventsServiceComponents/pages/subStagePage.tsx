
import { useParams } from "react-router-dom";
import SubStageMainPage  from "./subStageMainPage";

function SubStagePage() {
    const { id } = useParams<{ id: string }>();

    return (
        <SubStageMainPage eventId={id}></SubStageMainPage>
    );
}

export default SubStagePage;
