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
    classnumber: string
    disability: string
}

export interface User {
    firstname: string
    surname: string
    patronymic: string
    phone_number: string
    birthdate: string
    email: string
    gender: number
    role: number
}

export interface UserParticipant {
    User: User
    school: string
    disability: number
    classnumber: number
    citezenship: number
}