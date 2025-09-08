<?php

namespace Database\Seeders;

use App\Components\Dictionaries\AttendanceDictionary;
use App\Models\Attendance;
use App\Models\Task;
use App\Models\TaskAttendance;
use Illuminate\Database\Seeder;

class TaskAttendanceSeeder extends Seeder
{
    /**
     * Run the database seeds.
     *
     * @return void
     */
    public function run()
    {
        foreach (Task::all() as $task) {
            foreach (Attendance::where(['status' => AttendanceDictionary::ATTENDANCE])->get() as $attendance) {
                $taskAttendance = new TaskAttendance();
                $taskAttendance->task_id = $task->id;
                $taskAttendance->attendance_id = $attendance->id;
                $taskAttendance->points = rand(0, $task->max_points);
                $taskAttendance->save();
            }
        }
    }
}
