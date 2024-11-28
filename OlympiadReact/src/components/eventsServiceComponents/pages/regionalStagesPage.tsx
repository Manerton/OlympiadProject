import { REGIONAL_STAGE } from "../../../types/event"
import { RoleProvider } from "../../RoleContext"
import BaseEventPage from "./baseEventPage"

function RegionalStagesPage() {
    return (
    <RoleProvider>

        <BaseEventPage type={REGIONAL_STAGE}  pageName="Региональные этапы" />
    </RoleProvider>

    )
}

export default RegionalStagesPage

