export interface MyEvent {
    id?: string;
    name: string;

    dates: string[];
    profiles: string[];

    start_date: string;
    end_date: string;
    previous_event_id?: string;
    event_type: string;
    subject?: number;
    class_number?: number;
    additional_info?: string;
    events?: MyEvent[];
}

export interface MainEvent {
    id: string
    name: string

    dates: string[];
    profile: string;

    start_date: string
    end_date: string
    previous_event_id: string | null   // может быть null
    subject: number
    class_number: number
    additional_info: string
    event_type: string                 // есть в API
    events: MainEvent[] | null         // если внутри может быть список событий
}

export interface ApplicationEvent {
    MainEvent: MainEvent
    id: string
    status: number
    class_participation: number
}