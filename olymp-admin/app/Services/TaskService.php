<?php

namespace App\Services;

use App\Models\Task;
use App\Repositories\AttendanceRepository;
use App\Repositories\TaskAttendanceRepository;
use App\Repositories\TaskRepository;

class TaskService
{
    private TaskRepository $taskRepository;
    private TaskAttendanceRepository $taskAttendanceRepository;
    private AttendanceRepository $attendanceRepository;
    private ApplicationService $applicationService;
    public function __construct(
        TaskRepository $taskRepository,
        TaskAttendanceRepository $taskAttendanceRepository,
        AttendanceRepository $attendanceRepository,
        ApplicationService $applicationService
    )
    {
        $this->taskRepository = $taskRepository;
        $this->taskAttendanceRepository = $taskAttendanceRepository;
        $this->attendanceRepository = $attendanceRepository;
        $this->applicationService = $applicationService;
    }
    public function createTasks($data, $eventId){
        $applications = $this->applicationService->findByEventId($eventId);
        $applications = $this->applicationService->confirmedApplications($applications);
        $attendances = $this->attendanceRepository->getAttendances(array_column($applications, 'id'));
        foreach ($data['number'] as $key => $number) {
            $taskId = $this->taskRepository->create($eventId, $number, $data['point'][$key]);
            foreach ($attendances as $attendance) {
                $this->taskAttendanceRepository->create($attendance->id, $taskId, 0);
            }

        }
    }
    public function delete(Task $task)
    {
        $taskAttendances = $this->taskAttendanceRepository->getByTaskId($task->id);
        foreach ($taskAttendances as $taskAttendance) {
            $this->taskAttendanceRepository->delete($taskAttendance);
        }
        $this->taskRepository->delete($task);
    }
}
