<?php

namespace App\Repositories;

use App\Models\TaskAttendance;

class TaskAttendanceRepository
{
    private TaskRepository $taskRepository;
    public function __construct(
        TaskRepository $taskRepository
    )
    {
        $this->taskRepository = $taskRepository;
    }

    public function get($id){
        return TaskAttendance::find($id);
    }
    public function getAll(){
        return TaskAttendance::all();
    }
    public function getByTaskId($taskId){
        return TaskAttendance::where(['task_id' => $taskId])->get();
    }
    public function getByAttendanceId($attendanceId)
    {
        return TaskAttendance::where(['attendance_id' => $attendanceId])->get();
    }
    public function getByEventId($eventId){
        $tasks = $this->taskRepository->getByEventId($eventId);
        $taskAttendances = TaskAttendance::where(['task_id' => array_column($tasks->toArray(), 'id')])->get();
        return $taskAttendances;
    }
    public function getByAttendancesId($attendancesId)
    {
        return TaskAttendance::whereIn('attendance_id', $attendancesId)->get();
    }
    public function save(TaskAttendance $taskAttendance){
        return $taskAttendance->save();
    }
    public function delete(TaskAttendance $taskAttendance){
        return $taskAttendance->delete();
    }
    public function create($attendanceId, $taskId, $points)
    {
        $taskAttendance = new TaskAttendance();
        $taskAttendance->attendance_id = $attendanceId;
        $taskAttendance->task_id = $taskId;
        $taskAttendance->points = $points;
        $this->save($taskAttendance);
        return $taskAttendance->id;
    }
}
