<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class EventJury extends Model
{
    use HasFactory;
    public $event_id;
    public $user_id;
    public $eventAPI;
    public $userAPI;
}
