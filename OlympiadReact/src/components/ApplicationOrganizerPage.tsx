import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
interface Application {
  applicationID: number;
  userID: number;
  status: boolean | null; // true = одобрено, false = отклонено, null = не обработано
}

const ApplicationOrganizerPage: React.FC = () => {
  const { id } = useParams<{ id: string }>(); // типизируем id как строка;
  const [applications, setApplications] = useState<Application[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const fetchApplications = async () => {
      try {
        const response = await fetch(`http://localhost:8082/applications/event/${id}`, {
          credentials: 'include',
        });

        if (!response.ok) throw new Error("Failed to fetch applications");

        const data = await response.json();
        setApplications(data.data || []);
      } catch (error) {
        console.error('Error fetching applications:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchApplications();
  }, [id]);

  const handleStatusChange = async (applicationID: number, newStatus: boolean | null) => {
    try {
      const response = await fetch(`http://localhost:8082/applications/${applicationID}`, {
        method: 'PUT',
        body: JSON.stringify({ status: newStatus }),
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
      });
  
      if (!response.ok) throw new Error("Failed to update application status");
  
      // Обновить статус заявки в списке
      setApplications(prevApps =>
        prevApps.map(app =>
          app.applicationID === applicationID ? { ...app, status: newStatus } : app
        )
      );
    } catch (error) {
      console.error('Error updating status:', error);
    }
  };
  
  const handleDeleteApplication = async (applicationID: number) => {
    try {
      const response = await fetch(`http://localhost:8082/applications/${applicationID}`, {
        method: 'DELETE',
        credentials: 'include',
      });

      if (!response.ok) throw new Error("Failed to delete application");

      // Удалить заявку из списка
      setApplications(prevApps => prevApps.filter(app => app.applicationID !== applicationID));
    } catch (error) {
      console.error('Error deleting application:', error);
    }
  };

  if (loading) return <div>Loading...</div>;

  return (
    <div>
      <h1>Заявки на событие</h1>
      {loading ? (
        <div>Загрузка...</div>
      ) : (
        <table className="table">
          <thead>
            <tr>
              <th>ID пользователя</th>
              <th>Статус</th>
              <th>Действия</th>
            </tr>
          </thead>
          <tbody>
            {applications.map((app) => (
              <tr key={app.applicationID}>
                <td>{app.userID}</td>
                <td>
                  <select
                    value={app.status === null ? '' : app.status ? 'approved' : 'rejected'}
                    onChange={(e) =>
                      handleStatusChange(app.applicationID, e.target.value === 'approved' ? true : e.target.value === 'rejected' ? false : null)
                    }
                  >
                    <option value="">Не рассмотрено</option>
                    <option value="approved">Одобрено</option>
                    <option value="rejected">Отклонено</option>
                  </select>
                </td>
                <td>
                  <button
                    className="btn btn-danger"
                    onClick={() => handleDeleteApplication(app.applicationID)}
                  >
                    Удалить
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
  
}
export default ApplicationOrganizerPage;
