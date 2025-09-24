const API_URL = import.meta.env.VITE_API_URL;

export const API_CONFIG = {
    EVENT: `${API_URL}/events/`,
    CHILD: `${API_URL}/events/child`,
    STAGES: `${API_URL}/events/stages`,
    REGIONAL: `${API_URL}/regional-stages`,
    AUTH: `${API_URL}/api/users`,
    JUREASSIGNMENTS: `${API_URL}/jury-assignments`,
    APPLICATION: `${API_URL}/applications`,
    USERSBYROLE:  `${API_URL}/users/by-role/`,
    JURYBYSTAGE: `${API_URL}/jury-names/`,
    CREATEJURY: `${API_URL}/jury-assignments/create`,
    CREATEDELETEMANYJURY: `${API_URL}/jury-assignments/`,
    DELETEJURY: `${API_URL}/jury-assignments/remove/`,
};

export const AUTH = {
    login: `${API_URL}/auth`,
    logout: `${API_URL}/logout`,
    refresh: `${API_URL}/refresh`,
    register: `${API_URL}/register`,
    forgotPassword: `${API_URL}/users/forgot-password/`
};

export const SCHOOLS = {
    all: `${API_URL}/schools`,
    byId: `${API_URL}/schools/`,
};

export const USER = {
    info: `${API_URL}/users/all-info/`,
    changePassword: `${API_URL}/users/change-password/`,
};

export const APPLICATION = {
    getByUser: `${API_URL}/ApplicationEvent/`,
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
    info: `${API_URL}/users/participant/all-info/`,
};

export const HOSTS = {
    OLYMP_ADMIN: "http://localhost:8083",
    OLYMP_NOTIFICATION: "http://localhost:8084",
};

export const NOTIFY = {
    sendCode: `${API_URL}/email/send-code`
}

export default { API_CONFIG, HOSTS };
