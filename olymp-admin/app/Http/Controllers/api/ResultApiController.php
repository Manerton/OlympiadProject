<?php

namespace App\Http\Controllers\api;

use App\Http\Controllers\Controller;
use App\Repositories\AttendanceRepository;

class ResultApiController extends Controller
{
    private AttendanceRepository $attendanceRepository;
    public function __construct(
        AttendanceRepository $attendanceRepository
    )
    {
        $this->attendanceRepository = $attendanceRepository;
    }
    public function resultByAttendance($id){
        $attendance = $this->attendanceRepository->get($id);
        return response()->json([
            'result' => $attendance->taskAttendances
        ]);
    }
    public function resultByUser($id){

        return response()->json([]);
    }
    public function resultByUserTypeEvent($userId, $type, $eventId){

    }
}
