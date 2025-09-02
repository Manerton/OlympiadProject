export interface Appeal {
    TaskType: number;
    TaskID: number;
    Reason: string;
    Status: number;
}

export interface CreateAppealRequest {
    userId: string;
    taskId: number;
    reason: string;
}