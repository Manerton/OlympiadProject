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
    gender: number
    school: string
    classnumber: number
}