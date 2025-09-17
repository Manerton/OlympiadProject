<?php

namespace App\Components;

class ApiHelper
{
    public const STATUS_OK = 200;
    public const STATUS_NO_CONTENT = 204;
    public const STATUS_BAD_REQUEST = 400;
    public const STATUS_UNAUTHORIZED = 401;
    public const STATUS_NOT_FOUND = 404;
    public const OLYMP_NOTIFICATION_TOKEN = 1234567890;
    public const API_GATEWAY_HOST = 'http://apigateway-service:6611';
    public const NOTIFICATION_HOST = 'http://olymp-admin:8083';
    public const AUTH_URL_API = self::API_GATEWAY_HOST . '/auth';
    public const USER_URL_API = self::API_GATEWAY_HOST . '/users';
    public const TOKEN_REVOKE_URL_API = self::API_GATEWAY_HOST . '/revoke-all';
    public const TOKEN_REFRESH_URL_API = self::API_GATEWAY_HOST . '/refresh';
    public const EVENT_URL_API = self::API_GATEWAY_HOST . '/events/class';
    public const EVENT_MODEL_URL_API = self::API_GATEWAY_HOST . '/events';
    public const EVENT_JURY_URL_API = self::API_GATEWAY_HOST . '/jury-assignments/event';
    public const USER_COUNT_URL_API = self::API_GATEWAY_HOST . '/users/count';
    public const PARTICIPANT_URL_API = self::API_GATEWAY_HOST . '/participants';
    public const PARTICIPANT_COUNT_URL_API = self::API_GATEWAY_HOST . '/participants/count';
    public const PARTICIPANT_BY_USER_URL_API = self::API_GATEWAY_HOST . '/participants/byuser';
    public const SCHOOL_URL_API = self::API_GATEWAY_HOST . '/schools';
    public const SCHOOL_COUNT_URL_API = self::API_GATEWAY_HOST . '/schools/count';
    public const APPLICATION_URL_API = self::API_GATEWAY_HOST . '/applications';
    public const APPLICATION_EVENT_URL_API = self::API_GATEWAY_HOST . '/applications/event';
    public const APPLICATION_USER_URL_API = self::API_GATEWAY_HOST . '/applications/user';
    public const APPLICATION_COUNT_URL_API = self::API_GATEWAY_HOST . '/applications/count';
    public const LOGOUT_URL_API = self::API_GATEWAY_HOST . '/logout';
    public const SEND_CODE_URL_API = self::NOTIFICATION_HOST . '/index.php?r=email%2Fsend-code';
    public const SEND_MESSAGE_URL_API = self::NOTIFICATION_HOST . '/index.php?r=email%2Fsend-message';
    public const SUBJECTS_EVENT_URL_API = 'http://apigateway-service:8080/api/events/subjects';
    public static function prepareResponse($data)
    {
        return [
            'status' => 'success',
            'status_code' => 200,
            'data' => $data
        ];
    }
}
