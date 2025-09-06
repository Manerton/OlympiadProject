<?php

namespace App\Components\Dictionaries;

class NotificationTypeDictionary
{
    public const APPLICATION_REJECT = 'Ваша заявка была отклонена организаторами';
    public const APPLICATION_CONFIRM = 'Ваша заявка была подтверждена организаторами';
    public const APPEAL_REJECT = 'Ваша апелляция была подтверждена организаторами';
    public const APPEAL_CONFIRM = 'Ваша апелляция была подтверждена организаторами';
    public const RESULT_PUBLISH = 'Появились новые результаты';
    public const ONLINE_NOTIFICATION = 1;
}
