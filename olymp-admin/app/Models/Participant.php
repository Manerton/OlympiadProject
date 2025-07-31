<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class Participant extends Model
{
    use HasFactory;
    public User $userAPI;
    public School $schoolAPI;
    public $id;
    public $user_id;
    public $disability;
    public $citizenship;
    public $class;
    public $school_id;
    protected $fillable = [
        'user_id',
        'disability',
        'citizenship',
        'class',
        'school_id'
    ];
    protected $table = 'participant';
}
