<?php

namespace App\Repositories;

use App\Components\Dictionaries\AttendanceDictionary;
use App\Models\Attendance;

class AttendanceRepository
{
    public function get($id)
    {
        return Attendance::find($id);
    }
    public function getAll(){
        return Attendance::all();
    }
    public function delete(Attendance $attendance)
    {
        return $attendance->delete();
    }
    public function save(Attendance $attendance){
        return $attendance->save();
    }
    public function getAttendances($applicationsId)
    {
        return Attendance::whereIn('status', [AttendanceDictionary::ATTENDANCE, AttendanceDictionary::DISTANCE])
            ->whereIn('application_id', (array)$applicationsId)
            ->get();
    }
    public function getByApplicationId($applicationId)
    {
        return Attendance::where(['application_id' => $applicationId])->get();
    }
    public function getByApplicationsId($applicationId)
    {
        return Attendance::whereIn('application_id', $applicationId)->get();
    }
    public function create($applicationId){
        Attendance::create([
            'application_id' => $applicationId,
            'status' => AttendanceDictionary::NO_ATTENDANCE,
        ]);
    }
}
