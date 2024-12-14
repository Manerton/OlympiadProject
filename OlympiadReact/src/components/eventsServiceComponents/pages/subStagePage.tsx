
import { useParams } from "react-router-dom";

import { RoleProvider } from "../../RoleContext";
import SubStageMainPage  from "./subStageMainPage";

function SubStagePage() {

    const { id } = useParams<{ id: string }>();

    return (
        <RoleProvider>
            <SubStageMainPage eventId={id}></SubStageMainPage>
        </RoleProvider>
    );
}

export default SubStagePage;
