<?php

namespace App\Components\Dictionaries;

class AttendanceDictionary
{
    public const NO_ATTENDANCE = 1;
    public const ATTENDANCE = 2;
    public const DISTANCE = 3;
    public static function getList(){
        return [
            self::NO_ATTENDANCE => 'Неявка',
            self::ATTENDANCE => 'Явка',
            self::DISTANCE => 'Дистант'
        ];
    }
}
