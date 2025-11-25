import { useNavigate } from "react-router-dom";
import { useAuth } from "../../Helpers/AuthContext";
import { Profile} from "../../types/user";
import { UserRole } from "../../../dictionary/role";

import ChangePasswordBlock from "../Pages/Profile/ChangePasswordPage";
import { GetCitizenshipLabel } from "../../../dictionary/citizenship";

interface Props {
    profile: Profile
}

const PersonalInfo: React.FC<Props> = ({profile}) => {
    const { accessToken, user } = useAuth();
    const navigate = useNavigate();

    const actions = (
        <button className="btn btn-secondary" onClick={() => navigate("/profile")}>
            Назад
        </button>
    );

    return (
        <div className="d-flex flex-column">
            <div className="container">
                {/* Карточки профиля */}
                <div className="row g-3 mt-3">
                    <h5>Информация о вас</h5>
                    <div className="col-md-4">
                        <div className="card p-3 h-100 shadow-sm">
                            <h6 className="text-muted">Фамилия</h6>
                            <p className="fs-5 fw-semibold mb-0">{profile?.surname}</p>
                        </div>
                    </div>
                    <div className="col-md-4">
                        <div className="card p-3 h-100 shadow-sm">
                            <h6 className="text-muted">Имя</h6>
                            <p className="fs-5 fw-semibold mb-0">{profile?.firstname}</p>
                        </div>
                    </div>
                    <div className="col-md-4">
                        <div className="card p-3 h-100 shadow-sm">
                            <h6 className="text-muted">Отчество</h6>
                            <p className="fs-5 fw-semibold mb-0">{profile?.patronymic}</p>
                        </div>
                    </div>
                    <div className="col-md-4">
                        <div className="card p-3 h-100 shadow-sm">
                            <h6 className="text-muted">Телефон</h6>
                            <p className="fs-5 fw-semibold mb-0">{profile?.phone_number}</p>
                        </div>
                    </div>
                    <div className="col-md-4">
                        <div className="card p-3 h-100 shadow-sm">
                            <h6 className="text-muted">Дата рождения</h6>
                            <p className="fs-5 fw-semibold mb-0">{profile?.birthdate}</p>
                        </div>
                    </div>
                    <div className="col-md-4">
                        <div className="card p-3 h-100 shadow-sm">
                            <h6 className="text-muted">Пол</h6>
                            <p className="fs-5 fw-semibold mb-0">
                                {profile?.gender === 1 ? "Мужской" : profile?.gender === 2 ? "Женский" : ""}
                            </p>
                        </div>
                    </div>

                    {user?.role === UserRole.Participant && (
                        <>
                            <div className="col-md-6">
                                <div className="card p-3 h-100 shadow-sm">
                                    <h6 className="text-muted">Школа</h6>
                                    <p className="fs-5 fw-semibold mb-0">{profile?.school}</p>
                                </div>
                            </div>
                            <div className="col-md-3">
                                <div className="card p-3 h-100 shadow-sm">
                                    <h6 className="text-muted">Класс обучения</h6>
                                    <p className="fs-5 fw-semibold mb-0">{profile?.classnumber}</p>
                                </div>
                            </div>
                            <div className="col-md-3">
                                <div className="card p-3 h-100 shadow-sm">
                                    <h6 className="text-muted">Гражданство</h6>
                                    <p className="fs-5 fw-semibold mb-0">{GetCitizenshipLabel(profile?.citezenship)}</p>
                                </div>
                            </div>
                        </>
                    )}
                </div>

                {/* Блок смены пароля */}
                <div className="mt-5">
                    <ChangePasswordBlock />
                </div>
            </div>
        </div>
    );
};

export default PersonalInfo;
