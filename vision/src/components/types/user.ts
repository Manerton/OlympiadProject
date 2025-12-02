export interface UserAuth {
    id: string
    Email: string
    role: number
}


export interface RegisterForm {
    email: string
    password: string
    firstname: string
    surname: string
    patronymic: string
    phone_number: string
    birthdate: string
    gender: string
    school_id: string
    class_number: string
    disability: string
    citizenship: string
}

export interface User {
    id: string
    firstname: string
    surname: string
    patronymic: string
    phone_number: string
    birthdate: string
    email: string
    gender: number
    role: number
}

export interface Profile {
    user_id: string,
    participant_id: string,
    surname: string,
    firstname: string,
    patronymic: string,
    phone_number: string,
    birthdate: string,
    gender: number,
    school: string,
    classnumber: number,
    citezenship: number,
    email: string,
}

export interface UserParticipant {
    User: User
    participant_id: string
    school: string
    disability: number
    classnumber: number
    citezenship: number
}

export interface ChangePasswordForm {
    user_id: string
    old_password: string
    new_password: string
}

export interface ForgotPasswordForm {
    mail: string,
    code: string,
    password: string
}