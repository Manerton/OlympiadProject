<?php

namespace App\Components;

class RabbitMQHelper
{
    public const CREATE = 'create';
    public const UPDATE = 'update';
    public const DELETE = 'delete';
    public const QUEUE_NAME = 'olymp_admin';
    public const AUTH_QUEUE_NAME = 'olymp_auth';
    public const NOTIFICATION_QUEUE_NAME = 'olymp_notification';
    public const APPLICATION_QUEUE_NAME = 'olymp_application';
    public const EVENT_QUEUE_NAME = 'olymp_event';
    public const USER_TABLE = 'user';
    public const PARTICIPANT_TABLE = 'participant';
    public const SCHOOL_TABLE = 'school';
    public const APPLICATION_TABLE = 'applications';
}
