<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class Appeal extends Model
{
    use HasFactory;

    protected $table = 'appeal';
    protected $fillable = [
        'user_id', 'task_id', 'reason', 'status'
    ];
    public function task()
    {
        return $this->hasOne(Task::class, 'id', 'task_id');
    }
    public function user(){
        return $this->hasOne(User::class, 'id', 'task_id');
    }
}
