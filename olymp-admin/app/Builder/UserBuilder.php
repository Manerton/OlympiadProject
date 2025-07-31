<?php

namespace App\Builder;

use App\Models\Participant;
use App\Models\User;
use DateTime;

class UserBuilder
{
    public function build($item)
    {
        $model = new User();
        $model->id = $item['id'];
        $model->email = $item['email'];
        $model->firstname = $item['firstname'];
        $model->surname = $item['surname'];
        $model->patronymic = $item['patronymic'];
        $model->phone_number = $item['phone_number'];
        $model->gender = $item['gender'];
        $model->role = $item['role'];
        $date = new DateTime($item['birthdate']);
        $model->birthdate = $date->format('Y-m-d');
        return $model;
    }
    public function buildParticipant(User $model, Participant $participant)
    {
        $model->participantAPI = $participant;
    }
}
