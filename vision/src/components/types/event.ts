export interface MyEvent {
    id?: string;
    name: string;
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
    start_date: string
    end_date: string
    previous_event_id: string;
    subject: number;
    class_number: number
    additional_info: string
}

export interface ApplicationEvent {
    MainEvent: MainEvent
    status: number
}