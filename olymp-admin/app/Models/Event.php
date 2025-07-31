<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class Event extends Model
{
    use HasFactory;
    public $id;
    public $name;
    public $start_date;
    public $end_date;
    public $event_type;
    public $class_number;
    public $previous_event_id;
    public $subject;
    public $additional_info;
    public $events;
    public function tasks(){
        return $this->hasMany(Task::class, 'event_id', 'id');
    }

}
