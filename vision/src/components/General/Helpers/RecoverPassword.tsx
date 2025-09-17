import React, { useState } from "react";
import { axiosSendCode } from "../../../requests/NotificationRequests";
import { useAuth } from "../../Helpers/AuthContext";
import { axiosSSOForgotPassword } from "../../../requests/SSORequests";
import { ForgotPasswordForm } from "../../types/user";


async function verifyResetCode(email: string, code: string, password: string) {
    console.log("Проверка кода:", code, "для email:", email);
    return new Promise((resolve, reject) => {
        setTimeout(() => {
            if (code === "1234") resolve(true);
            else reject(new Error("Неверный код"));
        }, 1000);
    });
}

const ForgotPasswordPage: React.FC = () => {

    const { accessToken } = useAuth()

    const [step, setStep] = useState<1 | 2>(1);
    const [email, setEmail] = useState("");
    const [code, setCode] = useState("");
    const [password, setPassword] = useState("");
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const handleSendCode = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);
        setLoading(true);
        try {
            await axiosSendCode(email)
            setStep(2);
        } catch (err) {
            setError("Ошибка при отправке кода");
        } finally {
            setLoading(false);
        }
    };

    const handleResetPassword = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);
        setLoading(true);
        try {
            const forgotForm: ForgotPasswordForm = { 
                email: email,
                code: code,
                password: password
            }

            await axiosSSOForgotPassword( forgotForm)
            alert("Код подтверждён ✅ (тут можно открыть форму смены пароля)");
        } catch (err) {
            setError("Неверный код, попробуйте снова");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="container mt-5" style={{ maxWidth: 400 }}>
            <h3 className="mb-4 text-center">Восстановление пароля</h3>

            {step === 1 && (
                <form onSubmit={handleSendCode}>
                    <div className="mb-3">
                        <label className="form-label">Email</label>
                        <input
                            type="email"
                            className="form-control"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            required
                        />
                    </div>
                    {error && <p className="text-danger">{error}</p>}
                    <button type="submit" className="btn btn-primary w-100" disabled={loading}>
                        {loading ? "Отправка..." : "Отправить код"}
                    </button>
                </form>
            )}

            {step === 2 && (
                <form onSubmit={handleResetPassword}>
                    <div className="mb-3">
                        <label className="form-label">Код подтверждения</label>
                        <input
                            type="text"
                            className="form-control"
                            value={code}
                            onChange={(e) => setCode(e.target.value)}
                            required
                        />
                    </div>

                    <div className="mb-3">
                        <label className="form-label">Новый пароль</label>
                        <input
                            type="password"
                            className="form-control"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            required
                        />
                    </div>

                    {error && <p className="text-danger">{error}</p>}
                    <button type="submit" className="btn btn-success w-100" disabled={loading}>
                        {loading ? "Обновление..." : "Сбросить пароль"}
                    </button>
                </form>

            )}
        </div>
    );
}

export default ForgotPasswordPage