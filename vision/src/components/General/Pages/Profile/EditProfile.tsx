import React from "react";
import { useAuth } from "../../../Helpers/AuthContext";
import { useNavigate } from "react-router-dom";
import UserInfoBlock from "../../../Common/UserInfo";

const EditProfile: React.FC = () => {

    const { user } = useAuth()

    const navigate = useNavigate()


    // Пример начальных значений, заменить на реальные из user
    const initialProfile = {
        surname: "Иванов",
        firstname: "Иван",
        patronymic: "Иванович",
        phone_number: "",
        birthdate: "",
        gender: 0,
        school: "",
        classnumber: 0,
        email: user?.Email || ""
    }

    const [profile, setProfile] = React.useState(initialProfile)

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
        const { name, value } = e.target
        setProfile(prev => ({ ...prev, [name]: value }))
    }

    const handleSave = () => {
        // TODO: реализовать сохранение профиля
        navigate("/profile")
    }

    const actions = (
        <>
            <button className="btn btn-primary" onClick={handleSave}>Сохранить</button>
            <button className="btn btn-secondary" onClick={() => navigate("/profile")}>Отмена</button>
        </>
    )

    return (
        <div className="d-flex flex-column min-vh-100">
            <div className="container">
                {/* Информация о пользователе */}
                <UserInfoBlock email={user?.Email} actions={actions} />
                <form className="mt-3" onSubmit={e => { e.preventDefault(); handleSave(); }}>
                    <div className="row g-3">
                        <div className="col-md-4">
                            <label className="form-label">Фамилия</label>
                            <input type="text" className="form-control" name="surname" value={profile.surname} disabled />
                        </div>
                        <div className="col-md-4">
                            <label className="form-label">Имя</label>
                            <input type="text" className="form-control" name="firstname" value={profile.firstname} disabled />
                        </div>
                        <div className="col-md-4">
                            <label className="form-label">Отчество</label>
                            <input type="text" className="form-control" name="patronymic" value={profile.patronymic} disabled />
                        </div>
                        <div className="col-md-4">
                            <label className="form-label">Телефон</label>
                            <input type="text" className="form-control" name="phone_number" value={profile.phone_number} onChange={handleChange} />
                        </div>
                        <div className="col-md-4">
                            <label className="form-label">Дата рождения</label>
                            <input type="date" className="form-control" name="birthdate" value={profile.birthdate} onChange={handleChange} />
                        </div>
                        <div className="col-md-4">
                            <label className="form-label">Пол</label>
                            <select className="form-select" name="gender" value={profile.gender} onChange={handleChange}>
                                <option value={0}>Мужской</option>
                                <option value={1}>Женский</option>
                            </select>
                        </div>
                        <div className="col-md-6">
                            <label className="form-label">Школа</label>
                            <input type="text" className="form-control" name="school" value={profile.school} onChange={handleChange} />
                        </div>
                        <div className="col-md-6">
                            <label className="form-label">Класс</label>
                            <input type="number" className="form-control" name="classnumber" value={profile.classnumber} onChange={handleChange} />
                        </div>
                    </div>
                </form>
            </div>
        </div>
    )
}

export default EditProfile