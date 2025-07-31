<?php

namespace App\Components\Dictionaries;

class RoleDictionary
{
    public const ADMIN = 1;
    public const PARTICIPANT = 2;
    public const JURY = 3;
    public const ORGANIZER   = 4;
    public static function getList(){
        return [
            self::ADMIN => 'Администратор',
            self::JURY => 'Член жюри',
            self::PARTICIPANT => 'Участник',
            self::ORGANIZER => 'Организатор'
        ];
    }
}
