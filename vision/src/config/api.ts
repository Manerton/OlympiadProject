const API_URL = import.meta.env.VITE_API_URL || "/api";

export const API_CONFIG = {
    ALLEVENTS: `${API_URL}/events/`,
    EVENT: `${API_URL}/events/`,
    CHILD: `${API_URL}/events/child`,
    AVAILABLE: `${API_URL}/available-event/`,
    STAGES: `${API_URL}/events/stages`,
    REGIONAL: `${API_URL}/regional-stages`,
    OLYMPIAD: `${API_URL}/olympiad-stages`,
    AUTH: `${API_URL}/api/users`,
    JUREASSIGNMENTS: `${API_URL}/jury-assignments`,
    APPLICATION: `${API_URL}/applications`,
    USERSBYROLE:  `${API_URL}/users/by-role/`,
    JURYBYSTAGE: `${API_URL}/jury-names/`,
    CREATEMANYJURY: `${API_URL}/jury-assignments/many`,
    DELETEMANYJURY: `${API_URL}/jury-assignments/delete/many`,
};

export const AUTH = {
    login: `${API_URL}/auth`,
    logout: `${API_URL}/logout`,
    refresh: `${API_URL}/refresh`,
    register: `${API_URL}/register`,
    forgotPassword: `${API_URL}/users/forgot-password/`,
    verifySMS: `${API_URL}/sms/verify-code/`,
    verifyEmail: `${API_URL}/auth/check-email/`,
    verifyPhone: `${API_URL}/auth/check-phone/`,
    district: `${API_URL}/districts/`,
    school: `${API_URL}/schools/district/`,
};

export const SCHOOLS = {
    all: `${API_URL}/schools/all`,
    byId: `${API_URL}/schools/`,
};

export const USER = {
    update: `${API_URL}/users/`,
    info: `${API_URL}/users/all-info/`,
    changePassword: `${API_URL}/users/change-password/`,
};

export const APPLICATION = {
    getByUser: `${API_URL}/ApplicationEvent/`,
    create: `${API_URL}/applications/create/`,
    delete: `${API_URL}/applications/`,
};

export const RESULT = {
    allEventsByUserId: `${API_URL}/history-event/`,
    allByEventIdUserId: `${API_URL}/result/result-by-event-user/`,
    allEventsWithAppealByUserId: `${API_URL}/events-appeal/`,
};

export const APPEAL = {
    create: `${API_URL}/appeal/store`,
    getAppealsByEventUser: `${API_URL}/appeal/appeal-event-user/`,
    get: "",
};

export const PARTICIPANT = {
    update: `${API_URL}/participants/`,
    info: `${API_URL}/users/participant/all-info/`,
};

export const HOSTS = {
    OLYMP_ADMIN: "http://localhost:8083",
    OLYMP_NOTIFICATION: "http://localhost:8084",
};

export const NOTIFY = {
    sendCode: `${API_URL}/email/send-code`,
    sendSMSCode: `${API_URL}/sms/send-code`,
}

export default { API_CONFIG, HOSTS };
