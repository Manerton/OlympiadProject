<?php

namespace App\Components\Dictionaries;

class SubjectDictionary
{
    public const HISTORY = 1;
    public const BIOLOGY = 2;
    public const GEOGRAPHY = 3;
    public const FOREIGN_LANGUAGE = 4;
    public const INFORMATICS = 5;
    public const LITERATURE = 6;
    public const MATH = 7;
    public const SOCIAL = 8;
    public const RUSSIAN = 9;
    public const PHYSICS = 10;
    public const CHEMISTRY = 11;
    public const SUBJECTS = [
        self::RUSSIAN => [
            'name' => 'Русский язык',
            'filepath' => 'templates/RUSSIAN.xlsx',
            'registerList' => [
                'name' => 'лист регистрации',
                'startCell' => ['A', 3]
            ],
            'auditoriumList' => [
                'name' => 'список по аудитории',
                'startCell' => ['A', 3]
            ],
            'attendanceList' => [
                'name' => 'явка',
                'startCell' => ['A', 3]
            ],
            'pointLists' => [
                [
                    'name' => '9 класс',
                    'codeCell' => ['B', 6],
                    'pointCell' => ['C', 6],
                    'ignoreColumns' => [],
                    'category' => [9]
                ],
                [
                    'name' => '10 класс',
                    'codeCell' => ['B', 6],
                    'pointCell' => ['C', 6],
                    'ignoreColumns' => [],
                    'category' => [10]
                ],
                [
                    'name' => '11 класс',
                    'codeCell' => ['B', 6],
                    'pointCell' => ['C', 6],
                    'ignoreColumns' => [],
                    'category' => [11]
                ],

            ],
            'ratingList' => [
                'name' => 'Предварительный рейтинг',
                'startCell' => ['A', 4]
            ],
            'formApplicationList' => [
                'name' => 'форма приложения к протоколу',
                'startCell' => ['A', 6]
            ],
            'formProtocolList' => [
                'name' => 'форма протокола обезличенная',
                'startCell' => ['A', 6]
            ],
            'formESUList' => [
                'name' => 'Форма ЭСУ',
                'startCell' => ['A', 5]
            ]
        ],
    ];
    public static function getList(){
        return [
            self::HISTORY => 'История',
            self::BIOLOGY => 'Биология',
            self::GEOGRAPHY => 'География',
            self::FOREIGN_LANGUAGE => 'Иностранный язык',
            self::INFORMATICS => 'Информатика и ИКТ',
            self::LITERATURE => 'Литература',
            self::MATH => 'Математика',
            self::SOCIAL => 'Обществознание',
            self::RUSSIAN => 'Русский язык',
            self::PHYSICS => 'Физика',
            self::CHEMISTRY => 'Химия',
        ];
    }
}
