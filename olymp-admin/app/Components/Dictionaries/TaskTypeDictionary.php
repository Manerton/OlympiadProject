<?php

namespace App\Components\Dictionaries;

class TaskTypeDictionary
{
    public const READING = 1;
    public const PRACTICE = 2;
    public const TEST = 3;
    public static function getList(){
        return [
            self::READING => 'Устное задание',
            self::PRACTICE => 'Практическое задание',
            self::TEST => 'Тестовое задание',
        ];
    }
}
