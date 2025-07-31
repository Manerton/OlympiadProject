import React, { useEffect, useState } from 'react';
import axios from 'axios';

const UserIndex = () => {
    const [users, setUsers] = useState([]);
    const [loading, setLoading] = useState(true);
    const [currentPage, setCurrentPage] = useState(1);
    const [totalUsers, setTotalUsers] = useState(0);
    const [perPage, setPerPage] = useState(10);

    const fetchUsers = (page = 1) => {
        setLoading(true);
        const token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQG1haWwucnUiLCJleHAiOjE3ODU0OTE5MzksImlkIjoiMGU2OTkxOTQtZjc4MS00NWE2LTg3Y2YtNTRhOTYyMzI1Y2YyIiwicm9sZSI6MX0.-bc6ZKSP6Lbv6rYO89ZV65iWVHxCrFlUDPjM81N1Dyc';
        
        axios.get(`http://olymp-admin-v2/api/user/index/${page}`, {
            headers: {
                'Authorization': token
            },
            withCredentials: true
        })
        .then(response => {
            setUsers(response.data.users);
            setTotalUsers(response.data.usersAmount);
            setPerPage(response.data.perPage);
            setCurrentPage(page); // Убедимся, что текущая страница синхронизирована
            setLoading(false);
        })
        .catch(error => {
            console.error("Ошибка при получении пользователей:", error);
            setLoading(false);
        });
    };
    
    useEffect(() => {
        fetchUsers(currentPage);
    }, [currentPage]);

    if (loading) return <p>Загрузка...</p>;

    const totalPages = Math.ceil(totalUsers / perPage);

    return (
        <div className="user-index">
            <h2>Список пользователей</h2>
            <button 
                onClick={() => alert('Форма добавления пользователя')} 
                className="btn btn-success"
            >
                Добавить пользователя
            </button>

                        <table className="table">
                <thead>
                    <tr>
                        <th>#</th>
                        <th>ФИО</th>
                        <th>Эл.почта</th>
                        <th>Номер телефона</th>
                        <th>Действия</th>
                    </tr>
                </thead>
                <tbody>
                    {users.map((user, index) => (
                        <tr key={user.id}>
                            <td>{(currentPage - 1) * perPage + index + 1}</td>
                            <td>{user.full_name}</td>
                            <td>{user.email}</td>
                            <td>{user.phone_number}</td>
                            <td>
                                <button className="btn btn-primary btn-sm">Просмотр</button>
                                <button className="btn btn-warning btn-sm">Редактировать</button>
                                <button 
                                    className="btn btn-danger btn-sm"
                                    onClick={() => alert(`Удалить пользователя ${user.id}`)}
                                >
                                    Удалить
                                </button>
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
            <div className="pagination">
                {Array.from({ length: totalPages }, (_, i) => i + 1).map(page => (
                    <button
                        key={page}
                        className={`btn btn-sm ${currentPage === page ? 'btn-primary' : 'btn-light'}`}
                        onClick={() => setCurrentPage(page)}
                        disabled={currentPage === page}
                    >
                        {page}
                    </button>
                ))}
            </div>
        </div>
    );
};

export default UserIndex;