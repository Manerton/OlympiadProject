export interface MyEvent {
    id?: number 
    name: string;
    start_date: Date;
    end_date: Date;
    previous_event_id?: number
    event_type: string;
    subject?: number;
    class_number?: number;
    additional_info?: string;
    events?: [MyEvent]
}

export const REGIONAL_STAGE = "REGIONAL_STAGE"
export const OLYMPIAD = "OLYMPIAD"
export const STAGE = "STAGE"
export const VIEW_WORKS = "VIEW_WORKS"
export const APPEAL = "APPEAL"
export const CLASS = "CLASS"


