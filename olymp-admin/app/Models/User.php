<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class User extends Model
{
    use HasFactory;
    public $id;
    public $email;
    public $password;
    public $firstname;
    public $surname;
    public $patronymic;
    public $phone_number;
    public $gender;
    public $role;
    public $birthdate;
    public $participantAPI;

    protected $table = 'user';
    protected $fillable = [
        'email',
        'password',
        'firstname',
        'surname',
        'phone_number',
        'gender',
        'role',
        'birthdate'
    ];
    public function getFullFio(){
        return $this->surname . ' ' . $this->firstname . ' ' . $this->patronymic;
    }

}
