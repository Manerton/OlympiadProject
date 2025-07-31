<?php

namespace App\Components;

use App\Components\Dictionaries\AttendanceDictionary;
use App\Components\Dictionaries\CountryDictionary;
use App\Components\Dictionaries\DisabilityDictionary;
use App\Components\Dictionaries\GenderDictionary;
use App\Components\Dictionaries\ReasonParticipantDictionary;
use App\Components\Dictionaries\SubjectDictionary;
use PhpOffice\PhpSpreadsheet\Cell\Coordinate;
use PhpOffice\PhpSpreadsheet\IOFactory;
use PhpOffice\PhpSpreadsheet\Writer\Xlsx;
use Symfony\Component\HttpFoundation\StreamedResponse;

class ExcelCreator
{
    public static function createList($subject, $data)
    {
        $config = SubjectDictionary::SUBJECTS[$subject];
        $templatePath = resource_path($config['filepath']);
        if (!file_exists($templatePath)) {
            throw new \Exception('Шаблонный файл не найден');
        }
        $spreadsheet = IOFactory::load($templatePath);
        //заполняем данными
        self::fillRegisterList($data, $spreadsheet, $config);
        self::fillAuditoriumList($data, $spreadsheet, $config);
        self::fillAttendanceList($data, $spreadsheet, $config);
        self::fillResultLists($data, $spreadsheet, $config);
        self::fillRatingList($data, $spreadsheet, $config);
        self::fillProtocolList($data, $spreadsheet, $config);
        self::fillFacelessProtocolList($data, $spreadsheet, $config);
        self::fillFormEsuList($data, $spreadsheet, $config);
        //
        $writer = new Xlsx($spreadsheet);
        $response = new StreamedResponse(
            function () use ($writer) {
                $writer->save('php://output');
            }
        );
        $response->headers->set('Content-Type', 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet');
        $response->headers->set('Content-Disposition', 'attachment;filename="' . $config['name'] . '_' . date('d-m-Y') . '.xlsx"');
        $response->headers->set('Cache-Control', 'max-age=0');
        return $response;
    }
    public static function fillRegisterList($data, $spreadsheet, $config){
        $sheet = $spreadsheet->getSheetByName($config['registerList']['name']);
        $counter = 0;
        foreach ($data as $event) {
            $column = Coordinate::columnIndexFromString($config['registerList']['startCell'][0]);
            foreach($event['participants'] as $participant){
                $sheet->setCellValueByColumnAndRow($column,  $counter + $config['registerList']['startCell'][1], $participant['person']->surname);
                $sheet->setCellValueByColumnAndRow($column + 1, $counter + $config['registerList']['startCell'][1], $participant['person']->firstname);
                $sheet->setCellValueByColumnAndRow($column + 2, $counter + $config['registerList']['startCell'][1], $participant['person']->patronymic);
                $sheet->setCellValueByColumnAndRow($column + 3, $counter + $config['registerList']['startCell'][1], $participant['person']->birthdate);
                $sheet->setCellValueByColumnAndRow($column + 4, $counter + $config['registerList']['startCell'][1], $event['event']->class_number);
                $sheet->setCellValueByColumnAndRow($column + 5, $counter + $config['registerList']['startCell'][1], $participant['person']->participantAPI->schoolAPI->name);
                $counter++;
            }
        }
    }
    public static function fillAuditoriumList($data, $spreadsheet, $config)
    {
        $sheet = $spreadsheet->getSheetByName($config['auditoriumList']['name']);
        $counter = 0;
        foreach ($data as $event) {
            $column = Coordinate::columnIndexFromString($config['auditoriumList']['startCell'][0]);
            foreach($event['participants'] as $participant){
                $sheet->setCellValueByColumnAndRow($column + 1, $counter + $config['auditoriumList']['startCell'][1], $participant['person']->surname);
                $sheet->setCellValueByColumnAndRow($column + 2, $counter + $config['auditoriumList']['startCell'][1], $participant['person']->firstname);
                $sheet->setCellValueByColumnAndRow($column + 3, $counter + $config['auditoriumList']['startCell'][1], $participant['person']->patronymic);
                $sheet->setCellValueByColumnAndRow($column + 4, $counter + $config['auditoriumList']['startCell'][1], $participant['person']->birthdate);
                $sheet->setCellValueByColumnAndRow($column + 5, $counter + $config['auditoriumList']['startCell'][1], $event['event']->class_number);
                $counter++;
            }
        }
    }
    public static function fillAttendanceList($data, $spreadsheet, $config){
        $sheet = $spreadsheet->getSheetByName($config['attendanceList']['name']);
        $counter = 0;
        foreach ($data as $event) {
            $column = Coordinate::columnIndexFromString($config['attendanceList']['startCell'][0]);
            foreach($event['participants'] as $participant){
                $sheet->setCellValueByColumnAndRow($column, $counter + $config['attendanceList']['startCell'][1], $participant['person']->surname);
                $sheet->setCellValueByColumnAndRow($column + 1, $counter + $config['attendanceList']['startCell'][1], $participant['person']->firstname);
                $sheet->setCellValueByColumnAndRow($column + 2, $counter + $config['attendanceList']['startCell'][1], $participant['person']->patronymic);
                $sheet->setCellValueByColumnAndRow($column + 3, $counter + $config['attendanceList']['startCell'][1], $participant['person']->birthdate);
                $sheet->setCellValueByColumnAndRow($column + 4, $counter + $config['attendanceList']['startCell'][1], $event['event']->class_number);
                $sheet->setCellValueByColumnAndRow($column + 5, $counter + $config['attendanceList']['startCell'][1], $participant['application']->code);
                $sheet->setCellValueByColumnAndRow($column + 6, $counter + $config['attendanceList']['startCell'][1], $participant['attendance']->status - 1);
                $counter++;
            }
        }
    }
    public static function fillResultLists($data, $spreadsheet, $config){

        foreach ($config['pointLists'] as $list){
            $counter = 0;
            $sheet = $spreadsheet->getSheetByName($list['name']);
            foreach ($data as $event){
                if (in_array($event['event']->class_number, $list['category'])){
                    $column = Coordinate::columnIndexFromString($list['codeCell'][0]);
                    foreach($event['participants'] as $participant){
                        if ($participant['attendance']->status == AttendanceDictionary::ATTENDANCE){
                            $sheet->setCellValueByColumnAndRow($column, $counter + $list['codeCell'][1], $participant['application']->code);
                            foreach ($participant['taskAttendances'] as $index => $taskAttendance) {
                                $sheet->setCellValueByColumnAndRow($column + $index + 1, $counter + $list['pointCell'][1], $taskAttendance->points);

                            }
                            $counter++;
                        }
                    }
                }
            }
        }
    }
    public static function fillRatingList($data, $spreadsheet, $config){
        $sheet = $spreadsheet->getSheetByName($config['ratingList']['name']);
        $counter = 0;
        foreach ($data as $event) {
            $column = Coordinate::columnIndexFromString($config['ratingList']['startCell'][0]);
            $sheet->setCellValueByColumnAndRow($column, $counter + $config['ratingList']['startCell'][1], $event['event']->class_number . ' ' . 'класс');
            $sheet->mergeCellsByColumnAndRow(
                $column,
                ($counter + $config['ratingList']['startCell'][1]),
                $column + 4,
                ($counter + $config['ratingList']['startCell'][1])
            );
            $sheet->getStyleByColumnAndRow(
                $column,
                $counter + $config['ratingList']['startCell'][1],
                $column + 4,
                $counter + $config['ratingList']['startCell'][1]
            )->applyFromArray([
                'alignment' => [
                    'horizontal' => \PhpOffice\PhpSpreadsheet\Style\Alignment::HORIZONTAL_CENTER,
                    'vertical' => \PhpOffice\PhpSpreadsheet\Style\Alignment::VERTICAL_CENTER,
                ],
                'font' => [
                    'bold' => true,
                ],
            ]);
            $counter++;
            foreach($event['participants'] as $participant){
                if ($participant['attendance']->status == AttendanceDictionary::ATTENDANCE) {
                    $sheet->setCellValueByColumnAndRow($column + 1, $counter + $config['ratingList']['startCell'][1], $participant['person']->getFullFio());
                    $sheet->setCellValueByColumnAndRow($column + 2, $counter + $config['ratingList']['startCell'][1], $participant['person']->participantAPI->schoolAPI->name);
                    $sheet->setCellValueByColumnAndRow($column + 3, $counter + $config['ratingList']['startCell'][1], $event['event']->class_number . ' ' . 'класс');
                    $sheet->setCellValueByColumnAndRow($column + 4, $counter + $config['ratingList']['startCell'][1], $participant['attendance']->getTotalScore());
                    $counter++;
                }
            }
        }
    }

    public static function fillProtocolList($data, $spreadsheet, $config){
        $sheet = $spreadsheet->getSheetByName($config['formApplicationList']['name']);
        $counter = 0;
        foreach ($data as $event) {
            $column = Coordinate::columnIndexFromString($config['formApplicationList']['startCell'][0]);
            $sheet->setCellValueByColumnAndRow($column, $counter + $config['formApplicationList']['startCell'][1], $event['event']->class_number . ' ' . 'класс');
            $sheet->mergeCellsByColumnAndRow(
                $column,
                ($counter + $config['formApplicationList']['startCell'][1]),
                $column + 6,
                ($counter + $config['formApplicationList']['startCell'][1])
            );
            $sheet->getStyleByColumnAndRow(
                $column,
                $counter + $config['formApplicationList']['startCell'][1],
                $column + 4,
                $counter + $config['formApplicationList']['startCell'][1]
            )->applyFromArray([
                'alignment' => [
                    'horizontal' => \PhpOffice\PhpSpreadsheet\Style\Alignment::HORIZONTAL_CENTER,
                    'vertical' => \PhpOffice\PhpSpreadsheet\Style\Alignment::VERTICAL_CENTER,
                ],
                'font' => [
                    'bold' => true,
                ],
            ]);
            $counter++;
            foreach($event['participants'] as $participant){
                if ($participant['attendance']->status == AttendanceDictionary::ATTENDANCE) {
                    $sheet->setCellValueByColumnAndRow($column + 1, $counter + $config['formApplicationList']['startCell'][1], $participant['person']->getFullFio());
                    $sheet->setCellValueByColumnAndRow($column + 2, $counter + $config['formApplicationList']['startCell'][1], $participant['application']->code);
                    $sheet->setCellValueByColumnAndRow($column + 3, $counter + $config['formApplicationList']['startCell'][1], $participant['person']->participantAPI->schoolAPI->name);
                    $sheet->setCellValueByColumnAndRow($column + 4, $counter + $config['formApplicationList']['startCell'][1], $event['event']->class_number . ' ' . 'класс');
                    $sheet->setCellValueByColumnAndRow($column + 5, $counter + $config['formApplicationList']['startCell'][1], $participant['attendance']->getTotalScore());
                    $sheet->setCellValueByColumnAndRow($column + 6, $counter + $config['formApplicationList']['startCell'][1], self::defineStatus($participant['attendance']->getTotalScore(), $event['prizeScore'], $event['winnerScore']));
                    $counter++;
                }
            }
        }
    }
    public static function fillFacelessProtocolList($data, $spreadsheet, $config){
        $sheet = $spreadsheet->getSheetByName($config['formProtocolList']['name']);
        $counter = 0;
        foreach ($data as $event) {
            $column = Coordinate::columnIndexFromString($config['formApplicationList']['startCell'][0]);
            $sheet->setCellValueByColumnAndRow($column, $counter + $config['formApplicationList']['startCell'][1], $event['event']->class_number . ' ' . 'класс');
            $sheet->mergeCellsByColumnAndRow(
                $column,
                ($counter + $config['formApplicationList']['startCell'][1]),
                $column + 5,
                ($counter + $config['formApplicationList']['startCell'][1])
            );
            $sheet->getStyleByColumnAndRow(
                $column,
                $counter + $config['formApplicationList']['startCell'][1],
                $column + 4,
                $counter + $config['formApplicationList']['startCell'][1]
            )->applyFromArray([
                'alignment' => [
                    'horizontal' => \PhpOffice\PhpSpreadsheet\Style\Alignment::HORIZONTAL_CENTER,
                    'vertical' => \PhpOffice\PhpSpreadsheet\Style\Alignment::VERTICAL_CENTER,
                ],
                'font' => [
                    'bold' => true,
                ],
            ]);
            $counter++;
            foreach($event['participants'] as $participant){
                if ($participant['attendance']->status == AttendanceDictionary::ATTENDANCE) {
                    $sheet->setCellValueByColumnAndRow($column + 1, $counter + $config['formApplicationList']['startCell'][1], $participant['application']->code);
                    $sheet->setCellValueByColumnAndRow($column + 2, $counter + $config['formApplicationList']['startCell'][1], $participant['person']->participantAPI->schoolAPI->name);
                    $sheet->setCellValueByColumnAndRow($column + 3, $counter + $config['formApplicationList']['startCell'][1], $event['event']->class_number . ' ' . 'класс');
                    $sheet->setCellValueByColumnAndRow($column + 4, $counter + $config['formApplicationList']['startCell'][1], $participant['attendance']->getTotalScore());
                    $sheet->setCellValueByColumnAndRow($column + 5, $counter + $config['formApplicationList']['startCell'][1], self::defineStatus($participant['attendance']->getTotalScore(), $event['prizeScore'], $event['winnerScore']));
                    $counter++;
                }
            }
        }
    }
    public static function fillFormEsuList($data, $spreadsheet, $config){
        $sheet = $spreadsheet->getSheetByName($config['formESUList']['name']);
        $counter = 0;
        foreach ($data as $event) {
            $column = Coordinate::columnIndexFromString($config['formESUList']['startCell'][0]);
            foreach($event['participants'] as $participant){
                if ($participant['attendance']->status == AttendanceDictionary::ATTENDANCE) {
                    $sheet->setCellValueByColumnAndRow($column, $counter + $config['formESUList']['startCell'][1], 'Астраханская область'); // переделать
                    $sheet->setCellValueByColumnAndRow($column + 1, $counter + $config['formESUList']['startCell'][1], $participant['person']->surname);
                    $sheet->setCellValueByColumnAndRow($column + 2, $counter + $config['formESUList']['startCell'][1], $participant['person']->firstname);
                    $sheet->setCellValueByColumnAndRow($column + 3, $counter + $config['formESUList']['startCell'][1], $participant['person']->patronymic);
                    $sheet->setCellValueByColumnAndRow($column + 4, $counter + $config['formESUList']['startCell'][1], GenderDictionary::getList()[$participant['person']->gender]);
                    $sheet->setCellValueByColumnAndRow($column + 5, $counter + $config['formESUList']['startCell'][1], $participant['person']->birthdate);
                    $sheet->setCellValueByColumnAndRow($column + 6, $counter + $config['formESUList']['startCell'][1], CountryDictionary::getList()[$participant['person']->participantAPI->citizenship]);
                    $sheet->setCellValueByColumnAndRow($column + 7, $counter + $config['formESUList']['startCell'][1], DisabilityDictionary::getList()[$participant['person']->participantAPI->disability]);
                    $sheet->setCellValueByColumnAndRow($column + 8, $counter + $config['formESUList']['startCell'][1], $participant['person']->participantAPI->schoolAPI->name);
                    $sheet->setCellValueByColumnAndRow($column + 9, $counter + $config['formESUList']['startCell'][1], $event['event']->class_number);
                    $sheet->setCellValueByColumnAndRow($column + 10, $counter + $config['formESUList']['startCell'][1], $participant['person']->participantAPI->class);
                    $sheet->setCellValueByColumnAndRow($column + 11, $counter + $config['formESUList']['startCell'][1], $participant['attendance']->getTotalScore());
                    $sheet->setCellValueByColumnAndRow($column + 12, $counter + $config['formESUList']['startCell'][1], self::defineStatus($participant['attendance']->getTotalScore(), $event['prizeScore'], $event['winnerScore']));
                    $sheet->setCellValueByColumnAndRow($column + 13, $counter + $config['formESUList']['startCell'][1], 'Прошлый год');
                    $sheet->setCellValueByColumnAndRow($column + 14, $counter + $config['formESUList']['startCell'][1], 'Муниципалитет');
                    $sheet->setCellValueByColumnAndRow($column + 15, $counter + $config['formESUList']['startCell'][1], ReasonParticipantDictionary::getList()[$participant['application']->reason]);
                    $counter++;
                }
            }
        }
    }
    public static function defineStatus($point, $prizeScore, $winnerScore)
    {
        if ($point > $winnerScore) {
            return 'Победитель';
        }
        else if ($point < $prizeScore) {
            return 'Участник';
        }
        else {
            return 'Призёр';
        }
    }
}
