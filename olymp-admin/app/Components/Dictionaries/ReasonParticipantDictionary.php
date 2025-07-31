<?php

namespace App\Components\Dictionaries;

class ReasonParticipantDictionary
{
    public const LAST_YEAR = 1;
    public const THIS_YEAR = 2;
    public static function getList(){
        return [
            self::LAST_YEAR => 'По итогам олимпиады прошлого года',
            self::THIS_YEAR => 'По итогам олимпиады текущего года'
        ];
    }
}
