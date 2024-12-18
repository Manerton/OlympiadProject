import React, { useState, useEffect } from "react";
import ApplicationCard from "./Application";
import { Application } from "../types/application";
import { useRole } from "./RoleContext";
import API_CONFIG from "../config/apiConfig";
import { MyEvent } from "../types/event";

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
    const { role, id, name, clearRoleData } = useRole();
   
    /* useEffect(() => {
      // Ensure 'id' is defined before making the request
      if (id) {
        fetch(`http://localhost:8082/applications/user/${id}`, {
          credentials: "include",
        })
          .then((response) => response.json())
          .then((data) => {
            console.log("Fetched data:", data);
            setApplications(data.data || []); // Устанавливаем преобразованные данные
            setIsLoading(false);
          })
          .catch((error) => {
            console.error("Error fetching data:", error);
            setIsLoading(false);
          });
      }
    }, [id]); // Add 'id' as a dependency to the useEffect */

    useEffect(() => {
      const fetchApplications = async () => {
        try {
          // Шаг 1: Получение заявок
          const appResponse = await fetch(
            `http://localhost:8082/applications/user/${id}`,
            { credentials: "include" }
          );
  
          if (!appResponse.ok) throw new Error("Failed to fetch applications");
  
          const appData = await appResponse.json();
  
          const applicationsData = appData.data || [];
  
          // Шаг 2: Получение данных о событиях для всех eventID
          const eventIDs = applicationsData.map((app: any) => app.eventID);
  
          const eventResponse = await fetch(
            `${API_CONFIG.EVENTS}/list`, // URL сервиса events
            {
              method: "POST",
              credentials: "include",
              headers: { "Content-Type": "application/json" },
              body: JSON.stringify({
                ids:  eventIDs , // Фильтруем события по ID
              }),
            }
          );
  
          if (!eventResponse.ok) throw new Error("Failed to fetch event details");
  
          const eventData = await eventResponse.json();
          const events = eventData.data || [];
          console.log("Events data:", events);
          // Шаг 3: Объединяем данные заявок и событий
          const mergedApplications = applicationsData.map((app: any) => {
            const event: MyEvent = events.find((e: MyEvent) => e.ID === app.eventID);

            console.log("Matched Event for App:", app.eventID, "->", event);

            return {
              applicationID: app.applicationID,
              userID: app.userID,
              eventID: app.eventID,
              eventName: event?.Name || "Не указано",
              eventDate: event?.StartDate || "",
              subject: event.Subject || "Не указано",
              status: app.status,
              submittedAt: app.submittedAt,
              updatedAt: app.updatedAt,
            };
          });
  
          setApplications(mergedApplications);
        } catch (error) {
          console.error("Error fetching applications:", error);
        } finally {
          setIsLoading(false);
        }
      };
  
      if (id) {
        fetchApplications();
      }
    }, [id]);

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