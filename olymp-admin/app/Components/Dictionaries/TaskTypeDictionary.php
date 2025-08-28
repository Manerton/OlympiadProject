<?php

namespace App\Components\Dictionaries;

class TaskTypeDictionary
{
    public const READING = 1;
    public const PRACTICE = 2;
    public const WRITING = 3;
    public const TEST = 4;
    public static function getList(){
        return [
            self::READING => 'Устное задание',
            self::PRACTICE => 'Практическое задание',
            self::WRITING => 'Письменная задание',
            self::TEST => 'Тестовое задание',
        ];
    }
}
