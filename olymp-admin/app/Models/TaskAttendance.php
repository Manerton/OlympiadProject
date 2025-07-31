<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class TaskAttendance extends Model
{
    use HasFactory;
    protected $table = 'task_attendance';
    protected $fillable = [
        'attendance_id', 'task_id', 'points'
    ];
    public function task()
    {
        return $this->hasOne(Task::class, 'id', 'task_id');
    }
    public function attendance(){
        return $this->hasOne(Attendance::class, 'id', 'attendance_id');
    }
}
