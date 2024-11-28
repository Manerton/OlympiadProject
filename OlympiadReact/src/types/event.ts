export interface MyEvent {
    ID?: number 
    Name: string;
    StartDate: string;
    EndDate: string;
    PreviousEventID?: number
    EventType: string;
    Subject?: string;
    AdditionalInfo?: string;
    Events?: [MyEvent]
}

export const REGIONAL_STAGE = "REGIONAL_STAGE"
export const OLYMPIAD = "OLYMPIAD"
export const STAGE = "STAGE"
