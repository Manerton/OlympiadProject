export interface Application {
    applicationID: number;
    userID: number;
    eventID: number;
    eventName:     string;
	  eventDate:     string; 
    subject:       string;
    status: boolean | null; // true = одобрено, false = отклонено, null = не обработано
    submittedAt: string; // Используем ISO строки для дат
    updatedAt: string;
  }
  