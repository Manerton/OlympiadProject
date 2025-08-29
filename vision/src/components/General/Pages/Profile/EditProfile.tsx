import React, { useEffect, useState } from "react";
import { useAuth } from "../../../Helpers/AuthContext";
import { useNavigate } from "react-router-dom";
import UserInfoBlock from "../../../Common/UserInfo";
import { axiosSSOUserInfo, axiosSSOUserParticipantInfo } from "../../../../requests/SSORequests";
import { UserRole } from "../../../../dictionary/role";
import type { User, UserParticipant } from "../../../types/user";

const EditProfile: React.FC = () => {
  const { accessToken, user } = useAuth();
  const navigate = useNavigate();

  const [profile, setProfile] = useState({
    surname: "",
    firstname: "",
    patronymic: "",
    phone_number: "",
    birthdate: "",
    gender: 0,
    school: "",
    classnumber: 0,
    email: user?.Email || "",
  });

  useEffect(() => {
    const fetchUser = async () => {
      try {
        if (!accessToken || !user) return;

        if (user.role === UserRole.Participant) {
          const data: UserParticipant = await axiosSSOUserParticipantInfo(accessToken, user.id);

          console.log(data);

          setProfile({
            surname: data.User.surname,
            firstname: data.User.firstname,
            patronymic: data.User.patronymic,
            phone_number: data.User.phone_number,
            birthdate: data.User.birthdate,
            gender: data.User.gender,
            school: data.school,
            classnumber: data.classnumber,
            email: data.User.email,
          });
        } else {
          const data: User = await axiosSSOUserInfo(accessToken, user.id);

          setProfile({
            surname: data.surname,
            firstname: data.firstname,
            patronymic: data.patronymic,
            phone_number: data.phone_number,
            birthdate: data.birthdate,
            gender: data.gender,
            school: "", // нет в User, оставляем пустым
            classnumber: 0, // нет в User
            email: data.email,
          });
        }
      } catch (err) {
        console.error("Ошибка загрузки профиля:", err);
      }
    };

    fetchUser();
  }, [accessToken, user]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setProfile((prev) => ({ ...prev, [name]: value }));
  };

  const handleSave = () => {
    // TODO: реализовать сохранение профиля
    navigate("/profile");
  };

  const actions = (
    <>
      {/* <button className="btn btn-primary" onClick={handleSave}>Сохранить</button> */}
      <button className="btn btn-secondary" onClick={() => navigate("/profile")}>
        Отмена
      </button>
    </>
  );

  return (
    <div className="d-flex flex-column min-vh-100">
      <div className="container">
        <UserInfoBlock email={user?.Email} actions={actions} />

        <form className="mt-3" onSubmit={(e) => { e.preventDefault(); handleSave(); }}>
          <div className="row g-3">
            {/* поля формы */}
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
                <option value={1}>Мужской</option>
                <option value={2}>Женский</option>
              </select>
            </div>
            <div className="col-md-6">
              <label className="form-label">Школа</label>
              <input type="text" className="form-control" name="school" value={profile.school} onChange={handleChange} disabled={user?.role !== UserRole.Participant} />
            </div>
            <div className="col-md-6">
              <label className="form-label">Класс</label>
              <input type="number" className="form-control" name="classnumber" value={profile.classnumber} onChange={handleChange} disabled={user?.role !== UserRole.Participant} />
            </div>
          </div>
        </form>
      </div>
    </div>
  );
};

export default EditProfile;
