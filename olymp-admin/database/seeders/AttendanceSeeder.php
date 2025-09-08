<?php

namespace Database\Seeders;

use App\Components\Dictionaries\AttendanceDictionary;
use App\Components\Dictionaries\StatusDictionary;
use App\Models\Attendance;
use Illuminate\Database\Seeder;

class AttendanceSeeder extends Seeder
{
    private const APPLICATIONS = ['1', '2', '3' ,'4' ,'5' ,'6' ,'7' ,'8' ,'9', '10', '11', '12'];
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
