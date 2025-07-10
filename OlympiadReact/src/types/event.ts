export interface MyEvent {
    ID?: number 
    Name: string;
    StartDate: Date;
    EndDate: Date;
    PreviousEventID?: number
    EventType: string;
    Subject?: string;
    ClassNumber?: number;
    AdditionalInfo?: string;
    Events?: [MyEvent]
}

export const REGIONAL_STAGE = "REGIONAL_STAGE"
export const OLYMPIAD = "OLYMPIAD"
export const STAGE = "STAGE"
export const VIEW_WORKS = "VIEW_WORKS"
export const APPEAL = "APPEAL"
export const CLASS = "CLASS"


