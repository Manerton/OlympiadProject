<?php

namespace App\Builder;

use App\Models\Participant;
use App\Models\School;
use App\Models\User;

class ParticipantBuilder
{
    public function build($item)
    {
        $model = new Participant();
        $model->id = $item['id'];
        $model->citizenship = $item['citizenship'];
        $model->disability = $item['disability'];
        $model->class = $item['class_number'];
        $model->user_id = $item['user_id'];
        $model->school_id = $item['school_id'];
        return $model;
    }
    public function buildSchool(Participant $participant, School $school)
    {
        $participant->schoolAPI = $school;
    }
    public function buildUser(Participant $participant, User $user){
        $participant->userAPI = $user;
    }
}
