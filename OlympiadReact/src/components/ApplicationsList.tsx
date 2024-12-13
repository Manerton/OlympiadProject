import React, { useState, useEffect } from "react";
import ApplicationCard from "./Application";
import { Application } from "../types/application";

/* const ApplicationsPage: React.FC = () => {
  const [applications, setApplications] = useState<Application[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchApplications = async () => {
      setIsLoading(true);
      try {
        const response = await fetch("/api/applications");
        const data = await response.json();
        setApplications(data);
      } catch (error) {
        console.error("Ошибка при загрузке заявок:", error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchApplications();
  }, []);

  return (
    <div className="container mt-4">
      <h1>Мои заявки</h1>
      {isLoading ? (
        <p>Загрузка...</p>
      ) : applications.length === 0 ? (
        <p>У вас нет заявок.</p>
      ) : (
        <div className="row">
          {applications.map((app) => (
            <div className="col-md-4" key={app.applicationID}>
              <ApplicationCard application={app} />
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default ApplicationsPage; */
/* const mockApplications: Application[] = [
    {
      applicationID: 1,
      userID: 123,
      eventID: 456,
      status: true,
      submittedAt: "2024-11-01T10:00:00Z",
      updatedAt: "2024-11-10T15:30:00Z",
    },
    {
      applicationID: 2,
      userID: 123,
      eventID: 789,
      status: false,
      submittedAt: "2024-11-05T14:00:00Z",
      updatedAt: "2024-11-12T09:15:00Z",
    },
    {
      applicationID: 3,
      userID: 123,
      eventID: 1011,
      status: null,
      submittedAt: "2024-11-07T18:45:00Z",
      updatedAt: "2024-11-14T11:00:00Z",
    },
  ]; */
  /* const convertKeysToCamelCase = (data: any[]) => {
    return data.map((item) => ({
      applicationID: item.application_id,
      userID: item.user_id,
      eventID: item.event_id,
      eventName: item.event_name, // ВРЕМЕННО
      eventLocation: item.event_location, // ВРЕМЕННО
      eventDate: item.event_date, // ВРЕМЕННО
      status: item.status,
      submittedAt: item.submitted_at,
      updatedAt: item.updated_at,
    }));
  }; */

  const ApplicationsPage: React.FC = () => {
    const [applications, setApplications] = useState<Application[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
  
    useEffect(() => {
      // Запрос к API с использованием fetch
      fetch("http://localhost:8082/applications", {credentials: "include"}) 
        .then((response) => {
          if (!response.ok) {
            throw new Error("Ошибка при загрузке данных.");
          }
          return response.json();
        })
        .then((data) => {
          console.log("Полученные данные:", data);
          //const convertedData = convertKeysToCamelCase(data.data || []);
          setApplications(data.data || []); // Устанавливаем преобразованные данные
          setIsLoading(false);
        })
        .catch((err) => {
          setError("Не удалось загрузить данные. Попробуйте позже.");
          setIsLoading(false);
        });
    }, []);
  
    return (
      <div className="container mt-4">
        <h1>Мои заявки</h1>
        {isLoading ? (
          <p>Загрузка...</p>
        ) : error ? (
          <p className="text-danger">{error}</p>
        ) : applications.length === 0 ? (
          <p>У вас нет заявок.</p>
        ) : (
          <div className="row">
            {applications.map((app) => (
              <div className="col-md-4" key={app.applicationID}>
                <ApplicationCard application={app} />
              </div>
            ))}
          </div>
        )}
      </div>
    );
  };
  
  export default ApplicationsPage;