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

export interface ApplicationEvent {
    id: string
    name: string
    start_date: string
    end_date: string
    previous_event_id: string;
    subject: number;
    class_number: number
    additional_info: string
    status: number
}