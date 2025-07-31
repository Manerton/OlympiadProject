<?php

namespace App\Components\Dictionaries;

class ApplicationStatusDictionary
{
    public const AWAITING = 1;
    public const APPROVED = 2;
    public const REJECTED = 3;
    public static function getList()
    {
        return [
            self::AWAITING => 'Ожидает подтверждения',
            self::APPROVED => 'Подтверждена',
            self::REJECTED => 'Отклонена'
        ];
    }
}
