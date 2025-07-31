<?php

namespace App\Services;

use App\Components\Dictionaries\AttendanceDictionary;
use App\Repositories\AttendanceRepository;
use App\Repositories\TaskAttendanceRepository;

class AttendanceService
{
    private ApplicationService $applicationService;
    private UserService $userService;
    private AttendanceRepository $attendanceRepository;
    private TaskAttendanceRepository $taskAttendanceRepository;
    public function __construct(
        ApplicationService $applicationService,
        UserService $userService,
        AttendanceRepository $attendanceRepository,
        TaskAttendanceRepository $taskAttendanceRepository
    )
    {
        $this->applicationService = $applicationService;
        $this->userService = $userService;
        $this->attendanceRepository = $attendanceRepository;
        $this->taskAttendanceRepository = $taskAttendanceRepository;
    }

    public function createTable($attendances)
    {
        $data = [];
        foreach ($attendances as $attendance) {
            $application = $this->applicationService->find($attendance->application_id);
            $person = $this->userService->find($application->user_id);
            $data[] = [
                'attendance' => $attendance,
                'person' => $person
            ];
        }
        return $data;
    }
    public function createExtraTable($attendances)
    {
        $data = [];
        foreach ($attendances as $attendance) {
            $taskAttendances = $this->taskAttendanceRepository->getByAttendancesId(array_column([$attendance], 'id'));
            $application = $this->applicationService->find($attendance->application_id);
            $person = $this->userService->find($application->user_id);
            $data[] = [
                'attendance' => $attendance,
                'application' => $application,
                'person' => $person,
                'taskAttendances' => $taskAttendances,
            ];
        }
        return $data;
    }
    public function applicationFilter($applications){
        $data = $this->attendanceRepository->getByApplicationsId(array_column($applications,'id'));
        return $data;
    }
    public function attendanceFilter($attendances)
    {
        return $attendances->filter(function ($attendance) {
            return $attendance->status == AttendanceDictionary::ATTENDANCE;
        });
    }

    public function delete()
    {

    }
}
