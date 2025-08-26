<?php

namespace App\Components;

class ApiHelper
{
    public const STATUS_OK = 200;
    public const STATUS_NO_CONTENT = 204;
    public const STATUS_BAD_REQUEST = 400;
    public const STATUS_UNAUTHORIZED = 401;
    public const STATUS_NOT_FOUND = 404;
    public const AUTH_URL_API = 'http://172.16.0.196:6611/auth';
    public const USER_URL_API = 'http://172.16.0.196:6611/users';
    public const TOKEN_REVOKE_URL_API = 'http://172.16.0.196:6611/revoke-all';
    public const TOKEN_REFRESH_URL_API = 'http://172.16.0.196:6611/refresh';
    public const EVENT_URL_API = 'http://172.16.0.196:6611/events/class';
    public const EVENT_MODEL_URL_API = 'http://172.16.0.196:6611/events';
    public const EVENT_JURY_URL_API = 'http://172.16.0.196:6611/jury-assignments/event';
    public const USER_COUNT_URL_API = 'http://172.16.0.196:6611/users/count';
    public const PARTICIPANT_URL_API = 'http://172.16.0.196:6611/participants';
    public const PARTICIPANT_COUNT_URL_API = 'http://172.16.0.196:6611/participants/count';
    public const PARTICIPANT_BY_USER_URL_API = 'http://172.16.0.196:6611/participants/byuser';
    public const SCHOOL_URL_API = 'http://172.16.0.196:6611/schools';
    public const SCHOOL_COUNT_URL_API = 'http://172.16.0.196:6611/schools/count';
    public const APPLICATION_URL_API = 'http://172.16.0.196:6611/applications';
    public const APPLICATION_EVENT_URL_API = 'http://172.16.0.196:6611/applications/event';
    public const APPLICATION_COUNT_URL_API = 'http://172.16.0.196:6611/applications/count';
    public const LOGOUT_URL_API = 'http://172.16.0.196:6611/logout';
    public const SUBJECTS_EVENT_URL_API = 'http://172.16.1.39:8080/api/events/subjects';
}
