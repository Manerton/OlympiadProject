<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class Task extends Model
{
    use HasFactory;
    protected $table = 'task';
    protected $fillable = [
        'event_id', 'max_points', 'number', 'type'
    ];
    public function taskAttendances(){
        return $this->hasMany(TaskAttendance::class);
    }
    public function appeals(){
        return $this->hasMany(Appeal::class);
    }
}
