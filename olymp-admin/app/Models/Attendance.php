<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class Attendance extends Model
{
    use HasFactory;
    protected $table = 'attendance';
    protected $fillable = [
        'id', 'application_id', 'status'
    ];
    /*public $id;
    public $application_id;
    public $status;*/
    public function changeStatus($status){
        $this->status = $status;
    }
    public function taskAttendances()
    {
        return $this->hasMany(TaskAttendance::class);
    }
    public function getTotalScore(){
        $points = 0;
        foreach ($this->taskAttendances as $taskAttendance) {
            $points = $points + $taskAttendance->points;
        }
        return $points;
    }
}
