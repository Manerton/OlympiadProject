const API_CONFIG = {
    EVENTS: "http://172.16.1.39:8080/api/events",
    REGIONAL: "http://172.16.1.39:8080/api/events/regional-stage",
    AUTH: "http://172.16.1.39:8181/api/users",
    JUREASSIGNMENTS: "http://172.16.1.39:8090/jury-assignments",
    APPLICATION: "http://172.16.1.39:8082/applications",
};

export const AUTH = {
    login: "http://172.16.0.196:6611/auth",
    logout: "http://172.16.0.196:6611/logout",
    refresh: "http://172.16.0.196:6611/refresh"
}

const LOCAL_CONFIG = {
    EVENTS: "http://localhost:8080/events",
    JUREASSIGNMENTS: "http://localhost:8090/jury-assignments",
    APPLICATION: "http://localhost:8082/applications",
}

export default API_CONFIG;