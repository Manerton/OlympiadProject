export interface Application {
    applicationID: number;
    userID: number;
    eventID: number;
    eventName:     string;    //ВРЕМЕННО
	  eventLocation: string;    //ВРЕМЕННО
	  eventDate:     string; //ВРЕМЕННО
    status: boolean | null; // true = одобрено, false = отклонено, null = не обработано
    submittedAt: string; // Используем ISO строки для дат
    updatedAt: string;
  }
  