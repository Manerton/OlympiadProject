<?php

namespace Database\Seeders;

use App\Components\Dictionaries\AttendanceDictionary;
use App\Components\Dictionaries\StatusDictionary;
use App\Models\Attendance;
use Illuminate\Database\Seeder;

class AttendanceSeeder extends Seeder
{
    private const APPLICATIONS = [
        '66666666-6666-6666-6666-666666666661',
        '66666666-6666-6666-6666-666666666662',
        '66666666-6666-6666-6666-666666666663',
        '66666666-6666-6666-6666-666666666664',
        '66666666-6666-6666-6666-666666666665',
        '66666666-6666-6666-6666-666666666666'
    ];
    /**
     * Run the database seeds.
     *
     * @return void
     */
    public function run()
    {
        foreach (self::APPLICATIONS as $applicationId) {
            $attendance = new Attendance();
            $attendance->application_id = $applicationId;
            $attendance->status = rand(AttendanceDictionary::NO_ATTENDANCE, AttendanceDictionary::ATTENDANCE);
            $attendance->save();
        }
    }
}
