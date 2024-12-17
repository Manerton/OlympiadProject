import { REGIONAL_STAGE } from "../../../types/event"
import { RoleProvider } from "../../RoleContext"
import BaseEventPage from "./baseEventPage"

function RegionalStagesPage() {
    return (
        <BaseEventPage type={REGIONAL_STAGE}  pageName="Региональные этапы" />
    )
}

export default RegionalStagesPage

